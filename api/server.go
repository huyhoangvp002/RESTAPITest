package api

import (
	db "RESTAPITest/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Querier
	router *gin.Engine
}

func NewServer(store db.Querier) *Server {
	server := &Server{
		store: store,
	}

	router := gin.Default()

	// Health check
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Category routes
	router.POST("/categories", server.createCategory)
	router.POST("/products", server.CreateProduct)
	router.POST("/discount", server.CreateDiscount)
	router.POST("/accounts", server.CreateAccount)
	router.POST("/login", server.Login)

	router.GET("/products/:id", server.GetProduct)
	router.GET("/products/categories", server.GetProductByCate)
	router.GET("/products", server.ListByMaxPrice)
	router.GET("/accounts", server.GetAccountByUsername)

	router.PATCH("/products", server.UdateProduct)
	// router.GET("/products", server.getProductByCateRequest)

	server.router = router
	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
