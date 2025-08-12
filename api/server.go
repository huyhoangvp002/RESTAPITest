package api

import (
	db "RESTAPITest/db/sqlc"
	"RESTAPITest/token"
	"RESTAPITest/util"
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"time"

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

	// Khởi tạo router và middleware (gồm rate limiter)
	server.setUpRouter()

	fmt.Println("[DEBUG] TimeOut:", server.config.AccessTokenDuration)
	fmt.Println("[DEBUG] Server:", server.config.ServerAddress)

	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	// Static + CORS
	router.Static("/static", "./static")
	router.Use(cors.Default())

	// === Rate limiter setup (per-IP token bucket) ===
	// Bạn có thể điều chỉnh các giá trị này từ server.config nếu muốn.
	capacity := 10.0                   // tối đa token
	refillRate := 1.0                  // token nạp lại mỗi giây
	tokensCost := 1.0                  // mỗi request tốn bao nhiêu token
	cleanupInterval := 5 * time.Minute // chu kỳ dọn dẹp bucket

	// map lưu bucket theo IP và mutex bảo vệ
	buckets := make(map[string]*token.TokenBucket)
	var bucketsMu sync.Mutex

	// Goroutine dọn dẹp bucket "không hoạt động"
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			bucketsMu.Lock()
			for ip, b := range buckets {
				// nếu bucket lâu không refill (không hoạt động) thì xóa
				if now.Sub(b.GetLastRefill()) > 2*cleanupInterval {
					delete(buckets, ip)
				}
			}
			bucketsMu.Unlock()
		}
	}()

	// actual middleware
	rateLimitMiddleware := func(c *gin.Context) {
		clientIP := c.ClientIP()

		// get or create bucket
		bucketsMu.Lock()
		b, ok := buckets[clientIP]
		if !ok {
			b = token.NewTokenBucket(capacity, refillRate)
			buckets[clientIP] = b
		}
		bucketsMu.Unlock()

		// Allow?
		if !b.Allow(tokensCost) {
			// thông báo Retry-After (đặt 1s ở đây, bạn có thể sửa)
			c.Header("Retry-After", "1")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many requests",
				"message": "Vượt quá giới hạn request. Vui lòng thử lại sau.",
			})
			c.Abort()
			return
		}

		c.Next()
	}

	// Áp dụng rate limiter cho toàn bộ router
	router.Use(rateLimitMiddleware)

	// === End rate limiter setup ===

	// Public routes
	router.POST("/login", server.Login)
	router.POST("/register", server.Register)

	router.GET("/products", server.SearchProductByName)
	router.GET("/account", server.AccountByID)
	router.GET("/auth/google/login", server.HandleGoogleLogin)
	router.GET("/auth/google/callback", server.HandleGoogleCallback)
	router.GET("/", server.redirect)

	router.POST("/api/webhook", server.WebHook)

	// Authenticated routes
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
	authRoutes.PATCH("/account_info", server.UpdateAccountInfo)

	authRoutes.DELETE("/account:id", server.DeleteAccount)
	authRoutes.DELETE("/cart/:id", server.DeleteCart)
	authRoutes.DELETE("/product/:id", roleMiddleware("admin", "dealer"), server.DeleteProduct)
	authRoutes.DELETE("/category/:id", roleMiddleware("admin"), server.DeleteCategories)

	server.router = router
}

// Start chạy server trên address truyền vào
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
