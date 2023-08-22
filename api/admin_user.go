package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/util"
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

type createAdminUserResponse struct {
	AdminUser adminUserResponse `json:"admin_user"`
	RoleList  []db.Role         `json:"role_list"`
}

func newAdminUserResponse(result db.AdminUserTxResult) createAdminUserResponse {
	adminUserRsp := adminUserResponse{
		ID:                result.AdminUser.ID,
		Account:           result.AdminUser.Account,
		FullName:          result.AdminUser.FullName,
		Status:            result.AdminUser.Status,
		PasswordChangedAt: result.AdminUser.PasswordChangedAt,
		CreatedAt:         result.AdminUser.CreatedAt,
	}

	return createAdminUserResponse{
		AdminUser: adminUserRsp,
		RoleList:  result.RoleList,
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
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newAdminUserResponse(result)

	ctx.JSON(http.StatusOK, rsp)
}
