package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
)

type createRoleRequest struct {
	Name          string  `json:"name" binding:"required"`
	PermissionsID []int64 `json:"permissions_id" binding:"required"`
}

func (server *Server) createRole(ctx *gin.Context) {
	var req createRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateRoleTxParams{
		Name:          req.Name,
		PermissionsID: req.PermissionsID,
	}

	result, err := server.store.CreateRoleTx(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

type updateRoleRequestQuery struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateRoleRequest struct {
	Name          string  `json:"name"`
	PermissionsID []int64 `json:"permissions_id"`
}

func (server *Server) updateRole(ctx *gin.Context) {
	var query updateRoleRequestQuery
	if err := ctx.ShouldBindUri(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateRoleRequest
	ctx.BindJSON(&req)

	arg := db.UpdateRoleTxParams{
		ID: query.ID,
	}

	if req.Name != "" {
		arg.Name = req.Name
	}

	if req.PermissionsID != nil {
		arg.PermissionsID = req.PermissionsID
	}

	result, err := server.store.UpdateRoleTx(context.Background(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

type listRoleRequest struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type RoleResponse struct {
	db.Role
	PermissionList []db.Permission `json:"permission_list"`
}

type listRoleResponse struct {
	Count int32          `json:"count"`
	Data  []RoleResponse `json:"data"`
}

func (server *Server) listRole(ctx *gin.Context) {
	var req listRoleRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListRolesParams{
		Limit:  req.PageSize,
		Offset: (req.Page - 1) * req.PageSize,
	}

	roles, err := server.store.ListRoles(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	counts, err := server.store.GetRolesCount(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var roleResponses []RoleResponse

	for _, role := range roles {
		permissions, err := server.store.ListPermissionsForRole(ctx, role.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		roleRsp := RoleResponse{
			Role:           role,
			PermissionList: permissions,
		}
		roleResponses = append(roleResponses, roleRsp)
	}

	listResponse := listRoleResponse{
		Count: int32(counts),
		Data:  roleResponses,
	}

	ctx.JSON(http.StatusOK, listResponse)
}

type getRoleRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getRole(ctx *gin.Context) {
	var req getRoleRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	role, err := server.store.GetRole(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	permissionList, err := server.store.ListPermissionsForRole(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := RoleResponse{
		Role:           role,
		PermissionList: permissionList,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type deleteRoleRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteRole(ctx *gin.Context) {
	var req deleteRoleRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.DeleteRoleTxParams{
		ID: req.ID,
	}

	result, err := server.store.DeleteRoleTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}
