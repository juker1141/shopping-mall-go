package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/admin/permission", server.createPermission)
	router.GET("/admin/permissions", server.listPermission)
	router.GET("/admin/permission/:id", server.getPermission)
	router.PATCH("/admin/permission/:id", server.updatePermission)

	router.POST("/admin/role", server.createRole)
	router.GET("/admin/roles", server.listRole)
	router.GET("/admin/role/:id", server.getRole)
	router.PATCH("/admin/role/:id", server.updateRole)
	router.DELETE("/admin/role/:id", server.deleteRole)

	router.POST("/admin/admin_user", server.createAdminUser)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
