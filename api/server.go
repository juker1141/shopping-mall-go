package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/token"
	"github.com/juker1141/shopping-mall-go/util"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("fullName", validFullName)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// 登入後台
	router.POST("/admin/login", server.loginAdminUser)

	// 權限
	router.POST("/admin/permission", server.createPermission)
	router.GET("/admin/permissions", server.listPermission)
	router.GET("/admin/permission/:id", server.getPermission)
	router.PATCH("/admin/permission/:id", server.updatePermission)

	// 角色
	router.POST("/admin/role", server.createRole)
	router.GET("/admin/roles", server.listRole)
	router.GET("/admin/role/:id", server.getRole)
	router.PATCH("/admin/role/:id", server.updateRole)
	router.DELETE("/admin/role/:id", server.deleteRole)

	// 使用者
	router.POST("/admin/admin_user", server.createAdminUser)

	server.router = router
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
