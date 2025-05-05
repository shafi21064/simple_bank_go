package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}

	router.GET("api/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/api/user", server.createUser)
	router.POST("/api/account", server.createAccount)
	router.GET("/api/account/:id", server.getAccount)
	router.GET("/api/accounts", server.listAccount)
	router.PUT("/api/account/update", server.updateAccount)
	router.DELETE("/api/account/:id", server.deleteAccount)

	router.POST("/api/transfer", server.createTransfer)

	server.router = router
	return server
}

func errorResponse(err error) *gin.H {
	return &gin.H{"error": err.Error()}
}
