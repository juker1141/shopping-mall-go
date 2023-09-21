package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/token"
	"github.com/juker1141/shopping-mall-go/val"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func permissionMiddleware(store db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, exists := ctx.Get(authorizationPayloadKey)
		if !exists {
			err := errors.New("can not get permission")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		typePayload, ok := payload.(*token.Payload)
		if !ok {
			err := errors.New("payload is not of type Payload")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		permissions_id, err := store.ListPermissionsIDByAccount(ctx, typePayload.Account)
		if err != nil {
			if err == db.ErrRecordNotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		requestedPath := ctx.FullPath()

		if userHasPermission(requestedPath, permissions_id) {
			ctx.Next()
		} else {
			err := errors.New("account does not have the relevant permissions")
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}
	}
}

func userHasPermission(requestedPath string, permissions_id []int64) bool {
	switch {
	case strings.Contains(requestedPath, "permission") ||
		strings.Contains(requestedPath, "role") ||
		strings.Contains(requestedPath, "manager"):
		if val.ContainsNumber(permissions_id, accountPermissionCode) {
			return true
		} else {
			return false
		}
	case
		strings.Contains(requestedPath, "member"):
		if val.ContainsNumber(permissions_id, memberPermissionCode) {
			return true
		} else {
			return false
		}
	case strings.Contains(requestedPath, "product"):
		if val.ContainsNumber(permissions_id, productPermissionCode) {
			return true
		} else {
			return false
		}
	case strings.Contains(requestedPath, "order"):
		if val.ContainsNumber(permissions_id, orderPermissionCode) {
			return true
		} else {
			return false
		}
	case strings.Contains(requestedPath, "coupon"):
		if val.ContainsNumber(permissions_id, couponPermissionCode) {
			return true
		} else {
			return false
		}
	case strings.Contains(requestedPath, "news"):
		if val.ContainsNumber(permissions_id, newsPermissionCode) {
			return true
		} else {
			return false
		}
	default:
		return false
	}
}
