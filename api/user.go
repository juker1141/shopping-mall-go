package api

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/juker1141/shopping-mall-go/val"
)

var defaultAvaterPath = filepath.Join("static", "avatar_images", "default_avatar.png")

type createUserRequest struct {
	Account         string                `form:"account" binding:"required,alphanum,min=8"`
	Email           string                `form:"email" binding:"required,email"`
	FullName        string                `form:"full_name" binding:"required,min=2,fullName"`
	Password        string                `form:"password" binding:"required,min=8"`
	GenderId        int32                 `form:"gender_id" binding:"required,number"`
	Cellphone       string                `form:"cellphone" binding:"required,cellPhone"`
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
		Cellphone:       req.Cellphone,
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

		if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") || strings.HasSuffix(file.Filename, ".gif") {
			err := fmt.Errorf("only non-GIF image files are allowed for upload")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		const maxSize = 5 << 20 // 5MB
		if file.Size > maxSize {
			err := fmt.Errorf("file size exceeds 5 MB")
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
		arg.AvatarUrl = defaultAvaterPath
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

type userResponseForAdmin struct {
	ID                int64       `json:"id"`
	Account           string      `json:"account"`
	Email             string      `json:"email"`
	FullName          string      `json:"full_name"`
	GenderID          pgtype.Int4 `json:"gender_id"`
	Cellphone         string      `json:"cellphone"`
	Address           string      `json:"address"`
	ShippingAddress   string      `json:"shipping_address"`
	PostCode          string      `json:"post_code"`
	HashedPassword    string      `json:"hashed_password"`
	Status            int32       `json:"status"`
	AvatarUrl         string      `json:"avatar_url"`
	PasswordChangedAt time.Time   `json:"password_changed_at"`
	CreatedAt         time.Time   `json:"created_at"`
}

func newUserResponseForAdmin(user db.User) userResponseForAdmin {
	return userResponseForAdmin{
		ID:                user.ID,
		Account:           user.Account,
		Email:             user.Email,
		FullName:          user.FullName,
		GenderID:          user.GenderID,
		Cellphone:         user.Cellphone,
		Address:           user.Address,
		ShippingAddress:   user.ShippingAddress,
		PostCode:          user.PostCode,
		Status:            user.Status,
		AvatarUrl:         user.AvatarUrl,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

type userRoutesUri struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateUserRequest struct {
	FullName        string                `form:"full_name"`
	OldPassword     string                `form:"old_password"`
	NewPassword     string                `form:"new_password"`
	Cellphone       string                `form:"cellphone"`
	Address         string                `form:"address"`
	ShippingAddress string                `form:"shipping_address"`
	PostCode        string                `form:"post_code"`
	Status          *int32                `form:"status"`
	AvatarFile      *multipart.FileHeader `form:"avatar_file"`
}

func (server *Server) updateUserByAdmin(ctx *gin.Context) {
	var uri userRoutesUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
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

	if req.Cellphone != "" {
		if err := val.ValidateCellPhone(req.Cellphone); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		arg.Cellphone = pgtype.Text{
			String: req.Cellphone,
			Valid:  true,
		}
	}

	if req.Address != "" {
		arg.Address = pgtype.Text{
			String: req.Address,
			Valid:  true,
		}
	}

	if req.ShippingAddress != "" {
		arg.ShippingAddress = pgtype.Text{
			String: req.ShippingAddress,
			Valid:  true,
		}
	}

	if req.PostCode != "" {
		arg.PostCode = pgtype.Text{
			String: req.PostCode,
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

	if req.OldPassword != "" && req.NewPassword != "" && req.OldPassword != req.NewPassword {
		user, err := server.store.GetUser(ctx, uri.ID)
		if err != nil {
			if err == db.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		err = util.CheckPassword(req.OldPassword, user.HashedPassword)
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
		arg.PasswordChangedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}

	if req.AvatarFile != nil {
		file, err := ctx.FormFile("avatar_file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") || strings.HasSuffix(file.Filename, ".gif") {
			err := fmt.Errorf("only non-GIF image files are allowed for upload")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		const maxSize = 5 << 20 // 5MB
		if file.Size > maxSize {
			err := fmt.Errorf("file size exceeds 5 MB")
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

		arg.AvatarUrl = pgtype.Text{
			String: targetPath,
			Valid:  true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getUserByAdmin(ctx *gin.Context) {
	var uri userRoutesUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, uri.ID)
	if err != nil {
		if err == db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponseForAdmin(user)

	ctx.JSON(http.StatusOK, rsp)
}

type listUsersQuery struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type listUsersResponse struct {
	Count int32                  `json:"count"`
	Data  []userResponseForAdmin `json:"data"`
}

func (server *Server) listUsersByAdmin(ctx *gin.Context) {
	var query listUsersQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  query.PageSize,
		Offset: (query.Page - 1) * query.PageSize,
	}

	users, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var data []userResponseForAdmin
	for _, user := range users {
		data = append(data, newUserResponseForAdmin(user))
	}

	counts, err := server.store.GetUsersCount(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := listUsersResponse{
		Count: int32(counts),
		Data:  data,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type deleteUserResult struct {
	Message string `json:"message"`
}

func (server *Server) deleteUserByAdmin(ctx *gin.Context) {
	var uri userRoutesUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteUser(ctx, uri.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	result := deleteUserResult{
		Message: "Delete member user success.",
	}

	ctx.JSON(http.StatusOK, result)
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
