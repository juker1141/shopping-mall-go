package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/util"
)

type createUserRequest struct {
	Account         string `json:"account" binding:"required,alphanum,min=8"`
	Email           string `json:"email" binding:"required,email"`
	FullName        string `json:"full_name" binding:"required,min=2,fullName"`
	Password        string `json:"password" binding:"required,min=8"`
	GenderId        int32  `json:"gender_id" binding:"required,number"`
	Phone           string `json:"phone" binding:"required,e164"`
	Address         string `json:"address" binding:"required"`
	ShippingAddress string `json:"shipping_address" binding:"required"`
	PostCode        string `json:"post_code" binding:"required"`
	AvatarUrl       string `json:"avatar_url"`
	Status          int32  `json:"status" binding:"required,number"`
}

type userResponse struct {
	ID        int64  `json:"id"`
	Account   string `json:"account"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Phone     string `json:"phone"`
	AvatarUrl string `json:"avatar_url"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Account:   user.Account,
		Email:     user.Email,
		FullName:  user.FullName,
		Phone:     user.Phone,
		AvatarUrl: user.AvatarUrl,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Account:  req.Account,
		Email:    req.Email,
		FullName: req.FullName,
		GenderID: pgtype.Int4{
			Int32: req.GenderId,
			Valid: true,
		},
		HashedPassword:  hashedPassword,
		Phone:           req.Phone,
		Address:         req.Address,
		ShippingAddress: req.ShippingAddress,
		PostCode:        req.PostCode,
		AvatarUrl:       "",
		Status:          req.Status,
	}

	if req.AvatarUrl != "" {
		arg.AvatarUrl = req.AvatarUrl
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.ForeignKeyViolation || errCode == db.UniqueViolation {
			ctx.JSON(http.StatusConflict, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.CreateCart(ctx, db.CreateCartParams{
		Owner: pgtype.Text{
			String: user.Account,
			Valid:  true,
		},
		TotalPrice: 0,
		FinalPrice: 0,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)

	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Account  string `json:"account" binding:"required,alphanum,min=8"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginUserResponse struct {
	AccessToken          string       `json:"access_token"`
	AccessTokenExpiresAt time.Time    `json:"access_token_expires_at"`
	User                 userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByAccount(ctx, req.Account)
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if user.Status == 0 {
		err := fmt.Errorf("user '%v' is in a disabled state", user.Account)
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Account, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
		User:                 newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}
