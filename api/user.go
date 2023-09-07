package api

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/util"
)

type createUserRequest struct {
	Account         string                `form:"account" binding:"required,alphanum,min=8"`
	Email           string                `form:"email" binding:"required,email"`
	FullName        string                `form:"full_name" binding:"required,min=2,fullName"`
	Password        string                `form:"password" binding:"required,min=8"`
	GenderId        int32                 `form:"gender_id" binding:"required,number"`
	Phone           string                `form:"phone" binding:"required,twPhone"`
	Address         string                `form:"address" binding:"required"`
	ShippingAddress string                `form:"shipping_address" binding:"required"`
	PostCode        string                `form:"post_code" binding:"required"`
	Status          int32                 `form:"status" binding:"required,number"`
	AvatarFile      *multipart.FileHeader `form:"avatar_file"`
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
		AvatarUrl: user.AvatarUrl,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
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
		Status:          req.Status,
	}

	if req.AvatarFile != nil {
		file, err := ctx.FormFile("avatar_file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		timestamp := time.Now().UnixNano()
		filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)
		targetPath := filepath.Join("static", "avatar_images", filename)

		err = ctx.SaveUploadedFile(req.AvatarFile, targetPath)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		arg.AvatarUrl = targetPath
	} else {
		arg.AvatarUrl = filepath.Join("static", "avatar_images", "default_avatar.png")
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
