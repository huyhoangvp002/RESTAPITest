package api

import (
	db "RESTAPITest/db/sqlc"
	"RESTAPITest/token"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createDiscountRequest struct {
	DiscountValue int32 `json:"discount_value" binding:"required,min=0,max=100"`
	ProductID     int64 `json:"product_id"`
}

func (server *Server) CreateDiscount(ctx *gin.Context) {
	var req createDiscountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	account_IDRaw, err := server.store.GetAccountIDByUsername(ctx, authPayload.Username)
	fmt.Println("[DEBUG]|Account id:", account_IDRaw)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	accountID := sql.NullInt32{
		Int32: int32(account_IDRaw),
		Valid: true,
	}

	productIDRaw, err := server.store.GetProdIDByAccountID(ctx, accountID)
	fmt.Println("[DEBUG]|Product id:", productIDRaw)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	productID := sql.NullInt32{
		Int32: int32(req.ProductID),
		Valid: true,
	}

	if productIDRaw != req.ProductID {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission denied: you do not own this product"})
	}

	arg := db.CreateDiscountParams{
		DiscountValue: req.DiscountValue,
		ProductID:     productID,
		AccountID:     accountID,
	}

	discount, err := server.store.CreateDiscount(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	price, err := server.store.GetPriceByID(ctx, req.ProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	discount_price := db.UpdateDiscountPriceParams{
		ID:            int64(req.ProductID),
		DiscountPrice: int32(Discount(int(req.DiscountValue), int(price))),
	}
	server.store.UpdateDiscountPrice(ctx, discount_price)
	ctx.JSON(http.StatusOK, discount)
}

func Discount(discount int, price int) int {
	return ((price * (100 - discount)) / 100)
}
