package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/token"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/juker1141/shopping-mall-go/val"
)

type createAdminUserRequest struct {
	Account  string `json:"account" binding:"required,alphanum,min=8"`
	FullName string `json:"fullName" binding:"required,min=2,fullName"`
	Status   int32  `json:"status" binding:"required,number"`
	Password string `json:"password" binding:"required,min=8"`
	RoleID   int64  `json:"roleId" binding:"required"`
}

type adminUserResponse struct {
	ID                int64     `json:"id"`
	Account           string    `json:"account"`
	FullName          string    `json:"fullName"`
	Status            int32     `json:"status"`
	RoleID            int64     `json:"roleId"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	CreatedAt         time.Time `json:"createdAt"`
}

type adminUserAPIResponse struct {
	adminUserResponse
	Role db.Role `json:"role"`
}

func newAdminUserResponse(adminUser db.AdminUser) adminUserResponse {
	return adminUserResponse{
		ID:                adminUser.ID,
		Account:           adminUser.Account,
		FullName:          adminUser.FullName,
		Status:            adminUser.Status,
		RoleID:            int64(adminUser.RoleID.Int32),
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

	role, err := server.store.GetRole(ctx, req.RoleID)
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateAdminUserParams{
		Account:        req.Account,
		FullName:       req.FullName,
		HashedPassword: hashedPassword,
		Status:         req.Status,
		RoleID: pgtype.Int4{
			Int32: int32(role.ID),
			Valid: true,
		},
	}

	adminUser, err := server.store.CreateAdminUser(context.Background(), arg)
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
		adminUserResponse: newAdminUserResponse(adminUser),
		Role:              role,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type listAdminUsersQuery struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"pageSize" binding:"required,min=5,max=10"`
}

type listAdminUsersResponse struct {
	Count int32                  `json:"count"`
	Data  []adminUserAPIResponse `json:"data"`
}

func (server *Server) listAdminUsers(ctx *gin.Context) {
	var query listAdminUsersQuery
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
		role, err := server.store.GetRole(ctx, int64(adminUser.RoleID.Int32))
		if err != nil {
			if err == db.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		data = append(data, adminUserAPIResponse{
			adminUserResponse: newAdminUserResponse(adminUser),
			Role:              role,
		})
	}

	counts, err := server.store.GetAdminUsersCount(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := listAdminUsersResponse{
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

	role, err := server.store.GetRole(ctx, int64(adminUser.RoleID.Int32))
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := adminUserAPIResponse{
		adminUserResponse: newAdminUserResponse(adminUser),
		Role:              role,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type updateAdminUserRequest struct {
	FullName    string `json:"fullName"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
	Status      *int32 `json:"status"`
	RoleID      int64  `json:"roleId"`
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

	arg := db.UpdateAdminUserParams{
		ID: uri.ID,
	}

	if req.FullName != "" {
		if err := val.ValidateFullName(req.FullName); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		arg.FullName = pgtype.Text{
			String: req.FullName,
			Valid:  true,
		}
	}

	if req.Status != nil {
		if val.IsValidStatus(int(*req.Status)) {
			arg.Status = pgtype.Int4{
				Int32: *req.Status,
				Valid: true,
			}
		} else {
			err := fmt.Errorf("invalid status input: %d", *req.Status)
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}

	if req.RoleID != 0 {
		arg.RoleID = pgtype.Int4{
			Int32: int32(req.RoleID),
			Valid: true,
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

		arg.HashedPassword = pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		}
	}

	adminUser, err := server.store.UpdateAdminUser(ctx, arg)
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	role, err := server.store.GetRole(ctx, int64(adminUser.RoleID.Int32))
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := adminUserAPIResponse{
		adminUserResponse: newAdminUserResponse(adminUser),
		Role:              role,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) deleteAdminUser(ctx *gin.Context) {
	var uri adminUserRoutesUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteAdminUser(ctx, uri.ID)
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := deleteResult{
		Message: "Delete adminUser success.",
	}

	ctx.JSON(http.StatusOK, rsp)
}

type getAdminUserInfoResponse struct {
	adminUserResponse
	PermissionList []db.Permission `json:"permissionList"`
}

func (server *Server) getAdminUserInfo(ctx *gin.Context) {
	payload, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		err := errors.New("can not get permission")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	typePayload, ok := payload.(*token.Payload)
	if !ok {
		err := errors.New("payload is not of type Payload")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	adminUser, err := server.store.GetAdminUserByAccount(ctx, typePayload.Account)
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	permissionList, err := server.store.ListPermissionsForAdminUser(ctx, adminUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := getAdminUserInfoResponse{
		adminUserResponse: newAdminUserResponse(adminUser),
		PermissionList:    permissionList,
	}
	ctx.JSON(http.StatusOK, rsp)
}

type loginAdminUserRequest struct {
	Account  string `json:"account" binding:"required,alphanum,min=8"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginAdminUserResponse struct {
	SessionID             uuid.UUID `json:"sessionId"`
	AccessToken           string    `json:"accessToken"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt"`
	RefreshToken          string    `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
	adminUserResponse
	PermissionList []db.Permission `json:"permissionList"`
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
		adminUserResponse:     newAdminUserResponse(adminUser),
		PermissionList:        permissionList,
	}

	ctx.JSON(http.StatusOK, rsp)
}
