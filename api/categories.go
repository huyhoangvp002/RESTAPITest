package api

import (
	db "RESTAPITest/db/sqlc"
	"RESTAPITest/token"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createCategoryReq struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (server *Server) createCategory(ctx *gin.Context) {
	var req createCategoryReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("JSON bind err")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	account_IDRaw, err := server.store.GetAccountIDByUsername(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	account_ID := sql.NullInt32{
		Int32: int32(account_IDRaw),
		Valid: true,
	}

	arg := db.CreateCategoryParams{
		Name:      req.Name,
		Type:      req.Type,
		AccountID: account_ID,
	}

	cate, err := server.store.CreateCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error"})
		return
	}
	ctx.JSON(http.StatusOK, cate)
}
