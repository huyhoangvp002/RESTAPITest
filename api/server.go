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

	server.router = router
	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
