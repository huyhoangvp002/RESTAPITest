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
	buckets map[compositeKey]*token.TokenBucket
	mutex   sync.Mutex
}

// compositeKey kết hợp IP + Method + Path
type compositeKey struct {
	IP     string
	Method string
	Path   string
}

// NewRateLimiter tạo rate limiter mới
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		buckets: make(map[compositeKey]*token.TokenBucket),
	}
}

// Middleware factory tạo middleware với cấu hình riêng
func (rl *RateLimiter) RateLimitMiddleware(capacity, tokenCost, refillRate float64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Tạo composite key duy nhất cho mỗi client + endpoint
		key := compositeKey{
			IP:     c.ClientIP(),
			Method: c.Request.Method,
			Path:   c.FullPath(), // Lấy path đã đăng ký trong router
		}

		// Lấy hoặc tạo bucket
		rl.mutex.Lock()
		bucket, exists := rl.buckets[key]
		if !exists {
			bucket = token.NewTokenBucket(capacity, refillRate)
			rl.buckets[key] = bucket
		}
		rl.mutex.Unlock()

		// Kiểm tra rate limit
		if !bucket.Allow(tokenCost) {
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

func (rl *RateLimiter) StartCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			rl.mutex.Lock()
			for key, bucket := range rl.buckets {
				// Sử dụng bucket.LastUsed() đã được triển khai
				if time.Since(bucket.LastUsed()) > 24*time.Hour {
					delete(rl.buckets, key)
				}
			}
			rl.mutex.Unlock()
		}
	}()
}
