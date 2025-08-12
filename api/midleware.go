package api

import (
	"RESTAPITest/token"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		token, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		payload, err := tokenMaker.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()

		// authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		// if len(authorizationHeader) == 0 {
		// 	err := errors.New("authorization header is not provided")
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		// 	return
		// }

		// field := strings.Fields(authorizationHeader)
		// if len(field) < 2 {
		// 	err := errors.New("invalid authorization header format")
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		// 	return
		// }
		// authorizationType := strings.ToLower(field[0])
		// if authorizationType != authorizationTypeBearer {
		// 	err := fmt.Errorf("unsupported authorization type %s", authorizationType)
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		// 	return
		// }
		// accessToken := field[1]
		// payLoad, err := tokenMaker.VerifyToken(accessToken)
		// if err != nil {
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		// 	return
		// }
		// ctx.Set(authorizationPayloadKey, payLoad)
		// ctx.Next()

	}
}

func roleMiddleware(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payLoadRaw, exists := ctx.Get(authorizationPayloadKey)
		if !exists {
			err := errors.New("authorization payload not found")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		payLoad, ok := payLoadRaw.(*token.Payload)
		if !ok {
			err := errors.New("invalid payload")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		for _, role := range roles {
			if payLoad.Role == role {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "permission denied"})

	}
}

type RateLimiter struct {
	buckets map[string]*token.TokenBucket
	mutex   sync.Mutex
	config  RateConfig
}

// RateConfig cấu hình rate limit
type RateConfig struct {
	Capacity   float64       // Số token tối đa
	RefillRate float64       // Token nạp lại mỗi giây
	TokensCost float64       // Số token tốn mỗi request
	Cleanup    time.Duration // Chu kỳ dọn dẹp bucket
}

// NewRateLimiter khởi tạo rate limiter
func NewRateLimiter(config RateConfig) *RateLimiter {
	rl := &RateLimiter{
		buckets: make(map[string]*token.TokenBucket),
		config:  config,
	}

	// Bắt đầu dọn dẹp định kỳ
	go rl.startCleanup()
	return rl
}

// Middleware gin
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		rl.mutex.Lock()
		bucket, exists := rl.buckets[clientIP]
		if !exists {
			bucket = token.NewTokenBucket(rl.config.Capacity, rl.config.RefillRate)
			rl.buckets[clientIP] = bucket
		}
		rl.mutex.Unlock()

		if !bucket.Allow(rl.config.TokensCost) {
			c.Header("Retry-After", time.Second.String())
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many requests",
				"message": "Vượt quá giới hạn request. Vui lòng thử lại sau.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Dọn dẹp bucket không dùng
func (rl *RateLimiter) startCleanup() {
	ticker := time.NewTicker(rl.config.Cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		for ip, bucket := range rl.buckets {
			// Sử dụng method đã export
			elapsed := time.Since(bucket.GetLastRefill())

			// Kiểm tra không hoạt động trong 2 chu kỳ dọn dẹp
			if elapsed > 2*rl.config.Cleanup {
				delete(rl.buckets, ip)
			}
		}
		rl.mutex.Unlock()
	}
}
