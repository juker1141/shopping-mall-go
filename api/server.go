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
	// renew token
	router.POST("/admin/tokens/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// 權限
	authRoutes.POST("/admin/permission", server.createPermission)
	authRoutes.GET("/admin/permissions", server.listPermission)
	authRoutes.GET("/admin/permission/:id", server.getPermission)
	authRoutes.PATCH("/admin/permission/:id", server.updatePermission)

	// 角色
	authRoutes.POST("/admin/role", server.createRole)
	authRoutes.GET("/admin/roles", server.listRole)
	authRoutes.GET("/admin/role/:id", server.getRole)
	authRoutes.PATCH("/admin/role/:id", server.updateRole)
	authRoutes.DELETE("/admin/role/:id", server.deleteRole)

	// 使用者
	authRoutes.POST("/admin/user", server.createAdminUser)
	authRoutes.GET("/admin/users", server.listAdminUser)
	authRoutes.GET("/admin/user/:id", server.getAdminUser)
	authRoutes.PATCH("/admin/user/:id", server.updateAdminUser)
	authRoutes.DELETE("/admin/user/:id", server.deleteAdminUser)

	server.router = router
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
