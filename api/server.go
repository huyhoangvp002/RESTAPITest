package api

import (
	db "RESTAPITest/db/sqlc"
	"RESTAPITest/token"
	"RESTAPITest/util"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store      db.Querier
	router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(config util.Config, store db.Querier) (*Server, error) {

	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setUpRouter()
	fmt.Println("[DEBUG] TimeOut:", server.config.AccessTokenDuration)
	fmt.Println("[DEBUG] Server:", server.config.ServerAddress)

	return server, nil
}

func (server *Server) setUpRouter() {

	router := gin.Default()
	router.POST("/login", server.Login)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/categories", server.createCategory)
	authRoutes.POST("/products", server.CreateProduct)
	authRoutes.POST("/discount", server.CreateDiscount)
	router.POST("/accounts", server.CreateAccount)

	authRoutes.GET("/products/:id", server.GetProduct)
	authRoutes.GET("/products/categories", server.GetProductByCate)
	authRoutes.GET("/products", server.ListByMaxPrice)
	authRoutes.GET("/accounts", server.GetAccountByUsername)
	authRoutes.GET("/products/all", server.ListProducts)
	authRoutes.GET("/listaccounts", server.ListAccounts)

	authRoutes.PATCH("/products", server.UdateProduct)
	// router.GET("/products", server.getProductByCateRequest)

	server.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
