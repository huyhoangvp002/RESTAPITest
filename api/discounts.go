package api

import (
	db "RESTAPITest/db/sqlc"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createDiscountRequest struct {
	DiscountValue int32 `json:"discount_value" binding:"required,min=0,max=100"`
	ProductID     int   `json:"product_id"`
}

func (server *Server) CreateDiscount(ctx *gin.Context) {
	var req createDiscountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	PID := sql.NullInt32{
		Int32: int32(req.ProductID),
		Valid: true,
	}
	arg := db.CreateDiscountParams{
		DiscountValue: req.DiscountValue,
		ProductID:     PID,
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
	price, err := server.store.GetPriceByID(ctx, int64(req.ProductID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
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
