package api

import (
	"RESTAPITest/token"
	"errors"
	"net/http"

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
