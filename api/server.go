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
		v.RegisterValidation("twPhone", validTaiwanPhone)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// 登入後台
	router.POST("/login", server.loginUser)
	router.POST("/admin/login", server.loginAdminUser)
	// renew token
	router.POST("/admin/tokens/renew_access", server.renewAccessToken)

	// user 註冊
	router.POST("/user", server.createUser)

	authRoutes := router.Group("/auth").Use(authMiddleware(server.tokenMaker))

	adminRoutes := router.Group("/admin").Use(authMiddleware(server.tokenMaker)).Use(permissionMiddleware(server.store))
	// 獲取圖片
	router.Static("/static", "./static")

	// 權限
	adminRoutes.POST("/permission", server.createPermission)
	adminRoutes.GET("/permissions", server.listPermission)
	adminRoutes.GET("/permission/:id", server.getPermission)
	adminRoutes.PATCH("/permission/:id", server.updatePermission)

	// 角色
	adminRoutes.POST("/role", server.createRole)
	adminRoutes.GET("/roles", server.listRole)
	adminRoutes.GET("/role/:id", server.getRole)
	adminRoutes.PATCH("/role/:id", server.updateRole)
	adminRoutes.DELETE("/role/:id", server.deleteRole)

	// 使用者
	router.POST("/admin/user", server.createAdminUser)
	adminRoutes.GET("/users", server.listAdminUser)
	adminRoutes.GET("/user/:id", server.getAdminUser)
	adminRoutes.PATCH("/user/:id", server.updateAdminUser)
	adminRoutes.DELETE("/user/:id", server.deleteAdminUser)

	// 前台使用者
	authRoutes.PATCH("/user/:id", server.updateUser)

	server.router = router
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	server.InitProject()
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
