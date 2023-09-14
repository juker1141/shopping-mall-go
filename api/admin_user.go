package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/juker1141/shopping-mall-go/val"
)

type createAdminUserRequest struct {
	Account  string  `json:"account" binding:"required,alphanum,min=8"`
	FullName string  `json:"full_name" binding:"required,min=2,fullName"`
	Status   int32   `json:"status" binding:"required,number"`
	Password string  `json:"password" binding:"required,min=8"`
	RolesID  []int64 `json:"roles_id" binding:"required"`
}

type adminUserResponse struct {
	ID                int64     `json:"id"`
	Account           string    `json:"account"`
	FullName          string    `json:"full_name"`
	Status            int32     `json:"status"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type adminUserAPIResponse struct {
	AdminUser adminUserResponse `json:"admin_user"`
	RoleList  []db.Role         `json:"role_list"`
}

func newAdminUserResponse(adminUser db.AdminUser) adminUserResponse {
	return adminUserResponse{
		ID:                adminUser.ID,
		Account:           adminUser.Account,
		FullName:          adminUser.FullName,
		Status:            adminUser.Status,
		PasswordChangedAt: adminUser.PasswordChangedAt,
		CreatedAt:         adminUser.CreatedAt,
	}
}

func (server *Server) createAdminUser(ctx *gin.Context) {
	var req createAdminUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(req.RolesID) <= 0 || req.RolesID == nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAdminUserTxParams{
		Account:        req.Account,
		FullName:       req.FullName,
		HashedPassword: hashedPassword,
		Status:         req.Status,
		RolesID:        req.RolesID,
	}

	result, err := server.store.CreateAdminUserTx(context.Background(), arg)
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.ForeignKeyViolation || errCode == db.UniqueViolation {
			ctx.JSON(http.StatusConflict, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := adminUserAPIResponse{
		AdminUser: newAdminUserResponse(result.AdminUser),
		RoleList:  result.RoleList,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type listAdminUserQuery struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type listAdminUserResponse struct {
	Count int32                  `json:"count"`
	Data  []adminUserAPIResponse `json:"data"`
}

func (server *Server) listAdminUsers(ctx *gin.Context) {
	var query listAdminUserQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAdminUsersParams{
		Limit:  query.PageSize,
		Offset: (query.Page - 1) * query.PageSize,
	}

	adminUsers, err := server.store.ListAdminUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var data []adminUserAPIResponse
	for _, adminUser := range adminUsers {
		roles, err := server.store.ListRolesForAdminUser(ctx, adminUser.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		data = append(data, adminUserAPIResponse{
			AdminUser: newAdminUserResponse(adminUser),
			RoleList:  roles,
		})
	}

	counts, err := server.store.GetAdminUsersCount(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := listAdminUserResponse{
		Count: int32(counts),
		Data:  data,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type adminUserRoutesUri struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAdminUser(ctx *gin.Context) {
	var uri adminUserRoutesUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	adminUser, err := server.store.GetAdminUser(ctx, uri.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	roles, err := server.store.ListRolesForAdminUser(ctx, adminUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := adminUserAPIResponse{
		AdminUser: newAdminUserResponse(adminUser),
		RoleList:  roles,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type updateAdminUserRequest struct {
	FullName    string   `json:"full_name"`
	OldPassword string   `json:"old_password"`
	NewPassword string   `json:"new_password"`
	Status      *int32   `json:"status"`
	RolesID     *[]int64 `json:"roles_id"`
}

func (server *Server) updateAdminUser(ctx *gin.Context) {
	var uri adminUserRoutesUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateAdminUserRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateAdminUserTxParams{
		ID: uri.ID,
	}

	if req.FullName != "" {
		if err := val.ValidateFullName(req.FullName); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		arg.FullName = req.FullName
	}

	if req.Status != nil {
		if val.IsValidStatus(int(*req.Status)) {
			arg.Status = req.Status
		} else {
			err := fmt.Errorf("invalid status input: %d", *req.Status)
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}

	if req.RolesID != nil {
		if len(*req.RolesID) > 0 {
			arg.RolesID = *req.RolesID
		} else {
			err := fmt.Errorf("at least one role is required")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}

	// 如果使用者有想要更改密碼
	if req.OldPassword != "" && req.NewPassword != "" && req.OldPassword != req.NewPassword {
		adminUser, err := server.store.GetAdminUser(ctx, uri.ID)
		if err != nil {
			if err == db.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		err = util.CheckPassword(req.OldPassword, adminUser.HashedPassword)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		hashedPassword, err := util.HashPassword(req.NewPassword)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		arg.HashedPassword = hashedPassword
	}

	result, err := server.store.UpdateAdminUserTx(ctx, arg)
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := adminUserAPIResponse{
		AdminUser: newAdminUserResponse(result.AdminUser),
		RoleList:  result.RoleList,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) deleteAdminUser(ctx *gin.Context) {
	var uri adminUserRoutesUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.DeleteAdminUserTxParams{
		ID: uri.ID,
	}

	result, err := server.store.DeleteAdminUserTx(ctx, arg)
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

type loginAdminUserRequest struct {
	Account  string `json:"account" binding:"required,alphanum,min=8"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginAdminUserResponse struct {
	SessionID             uuid.UUID         `json:"session_id"`
	AccessToken           string            `json:"access_token"`
	AccessTokenExpiresAt  time.Time         `json:"access_token_expires_at"`
	RefreshToken          string            `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time         `json:"refresh_token_expires_at"`
	AdminUser             adminUserResponse `json:"admin_user"`
	PermissionList        []db.Permission   `json:"permission_list"`
}

func (server *Server) loginAdminUser(ctx *gin.Context) {
	var req loginAdminUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	adminUser, err := server.store.GetAdminUserByAccount(ctx, req.Account)
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if adminUser.Status == 0 {
		err := fmt.Errorf("user '%v' is in a disabled state", adminUser.Account)
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, adminUser.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	permissionList, err := server.store.ListPermissionsForAdminUser(ctx, adminUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(adminUser.Account, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(adminUser.Account, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Account:      adminUser.Account,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginAdminUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		AdminUser:             newAdminUserResponse(adminUser),
		PermissionList:        permissionList,
	}

	ctx.JSON(http.StatusOK, rsp)
}
