package api

import (
	db "RESTAPITest/db/sqlc"
	"RESTAPITest/token"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addToCartRequest struct {
	Product_ID int `json:"product_id" binding:"required"`
	Value      int `json:"value" binding:"required,min=1"`
}

func (server *Server) AddToCart(ctx *gin.Context) {
	var req addToCartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	account_id_raw, err := server.store.GetAccountIDByUsername(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	account_id := sql.NullInt32{
		Int32: int32(account_id_raw),
		Valid: true,
	}

	product_id := sql.NullInt32{
		Int32: int32(req.Product_ID),
		Valid: true,
	}

	arg := db.CreateCartParams{
		Value:     int32(req.Value),
		AccountID: account_id.Int32,
		ProductID: product_id.Int32,
	}

	cart, err := server.store.CreateCart(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, cart)

}

type showCart struct {
	PageID   int `form:"page_id" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=10"`
}

func (server *Server) ShowCart(ctx *gin.Context) {
	var req showCart
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	account_id_raw, err := server.store.GetAccountIDByUsername(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.ListCartByAccountIDParams{
		AccountID: int32(account_id_raw),
		Limit:     int32(req.PageSize),
		Offset:    int32((req.PageID - 1) * req.PageSize),
	}

	cart, err := server.store.ListCartByAccountID(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, cart)

}
