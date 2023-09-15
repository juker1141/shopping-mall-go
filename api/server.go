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

const (
	accountPermissionCode = int64(1)
	productPermissionCode = int64(2)
	orderPermissionCode   = int64(3)
	couponPermissionCode  = int64(4)
	newsPermissionCode    = int64(5)
)

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
		v.RegisterValidation("cellPhone", validCellPhone)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	server.setupNoAuthRoutes(router)

	server.setupAuthRoutes(router)

	server.setupAdminRoutes(router)

	// 獲取靜態資源
	router.Static("/static", "./static")

	server.router = router
}

func (server *Server) setupNoAuthRoutes(router *gin.Engine) {
	// 使用者登入網頁
	router.POST("/login", server.loginUser)
	// 登入後台
	router.POST("/admin/login", server.loginAdminUser)
	// Renew Token
	router.POST("/admin/tokens/renew_access", server.renewAccessToken)
	// user 註冊
	router.POST("/user", server.createUser)
}

func (server *Server) setupAdminRoutes(router *gin.Engine) {
	adminRoutes := router.Group("/admin").Use(authMiddleware(server.tokenMaker)).Use(permissionMiddleware(server.store))

	// 權限
	adminRoutes.POST("/permission", server.createPermission)
	adminRoutes.GET("/permissions", server.listPermissions)
	adminRoutes.GET("/permission/:id", server.getPermission)
	adminRoutes.PATCH("/permission/:id", server.updatePermission)

	// 角色
	adminRoutes.POST("/role", server.createRole)
	adminRoutes.GET("/roles", server.listRoles)
	adminRoutes.GET("/role/:id", server.getRole)
	adminRoutes.PATCH("/role/:id", server.updateRole)
	adminRoutes.DELETE("/role/:id", server.deleteRole)

	// 使用者
	adminRoutes.POST("/manager-user", server.createAdminUser)
	adminRoutes.GET("/manager-users", server.listAdminUsers)
	adminRoutes.GET("/manager-user/:id", server.getAdminUser)
	adminRoutes.PATCH("/manager-user/:id", server.updateAdminUser)
	adminRoutes.DELETE("/manager-user/:id", server.deleteAdminUser)

	// 顧客資料
	adminRoutes.GET("/member-users", server.listUsersByAdmin)
	adminRoutes.GET("/member-user/:id", server.getUserByAdmin)
	adminRoutes.PATCH("/member-user/:id", server.updateUserByAdmin)
	adminRoutes.DELETE("/member-user/:id", server.deleteUserByAdmin)

	// 商品
	adminRoutes.POST("/product", server.createProduct)
	adminRoutes.GET("/products", server.listProducts)
	adminRoutes.GET("/product/:id", server.getProduct)
	adminRoutes.PATCH("/product/:id", server.updateProduct)
	adminRoutes.DELETE("/product/:id", server.deleteProduct)

	// 優惠卷
	adminRoutes.POST("/coupon", server.createCoupon)
	adminRoutes.PATCH("/coupon/:id", server.updateCoupon)
}

func (server *Server) setupAuthRoutes(router *gin.Engine) {
	// authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	// 前台使用者
	// authRoutes.PATCH("/user/:id", server.updateUser)
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	server.InitProject()
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
