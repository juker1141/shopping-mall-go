package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
)

type VerifyEmailQuery struct {
	EmailID    int64  `form:"emailId" binding:"required,min=1"`
	SecretCode string `form:"secretCode" binding:"required,min=32,max=128"`
}

type VerifyEmailResponse struct {
	IsVerified bool `json:"isVerified"`
}

func (server *Server) VerifyEmail(ctx *gin.Context) {
	var query VerifyEmailQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	txResult, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    query.EmailID,
		SecretCode: query.SecretCode,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := VerifyEmailResponse{
		IsVerified: txResult.User.IsEmailVerified,
	}

	ctx.JSON(http.StatusOK, rsp)
}
