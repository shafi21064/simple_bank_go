package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/shafi21064/simplebank/db/sqlc"
)

// Servers serves HTTP request for out bank services
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// start server
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

// Server for handel http resulets routes
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	router.GET("api/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/api/account", server.createAccount)
	router.GET("/api/account/:id", server.getAccount)

	server.router = router
	return server
}

func errorResponse(err error) *gin.H {
	return &gin.H{"error": err.Error()}
}
