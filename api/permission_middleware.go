package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juker1141/shopping-mall-go/val"
)

func checkPermission(ctx *gin.Context, requestedPath string, permissions_id []int64) {
	switch {
	case strings.Contains(requestedPath, "permission") ||
		strings.Contains(requestedPath, "role") ||
		strings.Contains(requestedPath, "/admin/user") ||
		strings.Contains(requestedPath, "user"):
		if val.ContainsNumber(permissions_id, 1) {
			ctx.Next()
		} else {
			err := errors.New("account does not have the relevant permissions")
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}
	case strings.Contains(requestedPath, "product"):
		if val.ContainsNumber(permissions_id, 2) {
			ctx.Next()
		} else {
			err := errors.New("account does not have the relevant permissions")
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}
	case strings.Contains(requestedPath, "order"):
		if val.ContainsNumber(permissions_id, 3) {
			ctx.Next()
		} else {
			err := errors.New("account does not have the relevant permissions")
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}
	case strings.Contains(requestedPath, "coupon"):
		if val.ContainsNumber(permissions_id, 4) {
			ctx.Next()
		} else {
			err := errors.New("account does not have the relevant permissions")
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}
	case strings.Contains(requestedPath, "news"):
		if val.ContainsNumber(permissions_id, 5) {
			ctx.Next()
		} else {
			err := errors.New("account does not have the relevant permissions")
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}
	default:
		ctx.Next()
	}
}
