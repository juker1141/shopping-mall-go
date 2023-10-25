package api

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/token"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/juker1141/shopping-mall-go/worker"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	router          *gin.Engine
	taskDistributor worker.TaskDistributor
}

const (
	accountPermissionCode = int64(1)
	memberPermissionCode  = int64(2)
	productPermissionCode = int64(3)
	orderPermissionCode   = int64(4)
	couponPermissionCode  = int64(5)
	newsPermissionCode    = int64(6)
)

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
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

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	router.Use(cors.New(corsConfig))

	server.setupNoAuthRoutes(router)

	server.setupAuthRoutes(router)

	server.setupAdminRoutes(router)

	// 獲取靜態資源
	router.Static("/static", "./static")

	server.router = router
}

func (server *Server) setupNoAuthRoutes(router *gin.Engine) {
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "pong",
		})
	})
	// 使用者登入網頁
	router.POST("/login", server.loginUser)
	// 登入後台
	router.POST("/admin/login", server.loginAdminUser)
	// Renew Token
	router.POST("/admin/tokens/renew_access", server.renewAccessToken)
	// 會員 註冊
	router.POST("/member_user", server.createUser)
	// 會員註冊驗證信 驗證
	router.GET("/member_user/verify_email", server.VerifyEmail)
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
	adminRoutes.GET("/roles/option", server.listRolesOption)
	adminRoutes.GET("/role/:id", server.getRole)
	adminRoutes.PATCH("/role/:id", server.updateRole)
	// adminRoutes.DELETE("/role/:id", server.deleteRole)

	// 管理者
	adminRoutes.POST("/manager_user", server.createAdminUser)
	adminRoutes.GET("/manager_users", server.listAdminUsers)
	adminRoutes.GET("/manager_user/:id", server.getAdminUser)
	adminRoutes.PATCH("/manager_user/:id", server.updateAdminUser)
	adminRoutes.DELETE("/manager_user/:id", server.deleteAdminUser)

	// 管理者刷新權限
	adminRoutes.GET("/manager_user/info", server.getAdminUserInfo)

	// 會員資料
	adminRoutes.GET("/member_users", server.listUsersByAdmin)
	adminRoutes.GET("/member_user/:id", server.getUserByAdmin)
	adminRoutes.PATCH("/member_user/:id", server.updateUserByAdmin)
	adminRoutes.DELETE("/member_user/:id", server.deleteUserByAdmin)

	// 商品
	adminRoutes.POST("/product", server.createProduct)
	adminRoutes.GET("/products", server.listProducts)
	adminRoutes.GET("/product/:id", server.getProduct)
	adminRoutes.PATCH("/product/:id", server.updateProduct)
	adminRoutes.DELETE("/product/:id", server.deleteProduct)

	// 訂單
	adminRoutes.POST("/order", server.createOrder)
	adminRoutes.GET("/orders", server.listOrders)
	adminRoutes.GET("/order/:id", server.getOrder)
	adminRoutes.PATCH("/order/:id", server.updateOrder)

	// 取得訂單狀態
	adminRoutes.GET("/order/statuses/option", server.listOrderStatusesOption)
	// 付款方式
	adminRoutes.GET("/order/pay_methods/option", server.listPayMethodsOption)

	// 優惠卷
	adminRoutes.POST("/coupon", server.createCoupon)
	adminRoutes.GET("/coupons", server.listCoupons)
	adminRoutes.GET("/coupon/:id", server.getCoupon)
	adminRoutes.PATCH("/coupon/:id", server.updateCoupon)
	adminRoutes.DELETE("/coupon/:id", server.deleteCoupon)

	// 檢查優惠卷是否有效
	adminRoutes.POST("/coupon/check", server.checkCoupon)
}

func (server *Server) setupAuthRoutes(router *gin.Engine) {
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	// 前台使用者
	// authRoutes.PATCH("/user/:id", server.updateUser)

	// 新增至購物車
	authRoutes.PATCH("/cart", server.updateCart)
	// 會員 下訂單
	authRoutes.POST("/order", server.createOrder)

	// 取得訂單狀態
	authRoutes.GET("/order/statuses/option", server.listOrderStatusesOption)
	// 付款方式
	authRoutes.GET("/order/pay_methods/option", server.listPayMethodsOption)
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	server.InitProject()
	return server.router.Run(address)
}

type deleteResult struct {
	Message string `json:"message"`
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
