package api

import (
	db "RESTAPITest/db/sqlc"
	"RESTAPITest/token"
	"RESTAPITest/util"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
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

	router.Static("/static", "./static")
	router.Use(cors.Default())

	router.POST("/login", server.Login)
	router.POST("/register", server.Register)

	router.GET("/products", server.SearchProductByName)
	router.GET("/auth/google/login", server.HandleGoogleLogin)
	router.GET("/auth/google/callback", server.HandleGoogleCallback)
	router.GET("/", server.redirect)

	router.POST("/api/webhook", server.WebHook)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/categories", roleMiddleware("admin"), server.createCategory)
	authRoutes.POST("/products", roleMiddleware("admin", "dealer"), server.CreateProduct)
	authRoutes.POST("/discount", roleMiddleware("admin", "dealer"), server.CreateDiscount)
	authRoutes.POST("/updaterole", roleMiddleware("admin"), server.UpdateRole)
	authRoutes.POST("/create_info", server.CreateAccountInfo)
	authRoutes.POST("/cart/add", roleMiddleware("admin", "customer"), server.AddToCart)
	authRoutes.POST("/orders", server.CreateOrder)
	authRoutes.POST("/shipments", server.CreateShipment)
	authRoutes.POST("/logout", server.LogOut)

	authRoutes.GET("/products/:id", server.GetProduct)
	authRoutes.GET("/products/categories", server.GetProductByCate)
	authRoutes.GET("/products_by_price", server.ListByMaxPrice)
	authRoutes.GET("/accounts", server.GetAccountByUsername)
	authRoutes.GET("/products/all", server.ListProducts)
	authRoutes.GET("/account/list", roleMiddleware("admin"), server.ListAccounts)
	authRoutes.GET("/cart", server.ShowCart)
	authRoutes.GET("/categories/all", server.ListCategories)
	authRoutes.GET("/profile", server.GetAccountInfo)

	authRoutes.PATCH("/products/:id", server.UdateProduct)
	authRoutes.PATCH("/cart/:id", server.UpdateProductInCart)
	authRoutes.PATCH("/account_info/:id", server.UpdateAccountInfo)

	authRoutes.DELETE("/account:id", server.DeleteAccount)
	authRoutes.DELETE("/cart/:id", server.DeleteCart)
	authRoutes.DELETE("/product/:id", roleMiddleware("admin", "dealer"), server.DeleteProduct)
	authRoutes.DELETE("/category/:id", roleMiddleware("admin"), server.DeleteCategories)

	server.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func statusCodeForError(err error) int {
	if err == sql.ErrNoRows {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
