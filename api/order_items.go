package api

import (
	db "RESTAPITest/db/sqlc"
	"RESTAPITest/token"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createOrderRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int64 `json:"quantity" binding:"required,min=1"`
}

func (server *Server) CreateOrder(ctx *gin.Context) {
	var req createOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload.Role != "buyer" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": "Permision Deny"})
		return
	}
	sellerID, err := server.store.GetAccountIDbyProductID(ctx, req.ProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	buyerID, err := server.store.GetAccountIDByUsername(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	price, err := server.store.GetDiscountPriceByID(ctx, req.ProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	oldStock, err := server.store.GetStockByID(ctx, req.ProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if oldStock == 0 {
		ctx.JSON(http.StatusConflict, gin.H{"error": "product is out of stock"})
		return
	} else if oldStock >= int32(req.Quantity) {
		arg := db.CreateOrderParams{
			BuyerID:    buyerID,
			SellerID:   sellerID,
			TotalPrice: int64(TotalPrice(int(price), int(req.Quantity))),
			Cod:        true,
			Status:     "pending",
		}

		order, err := server.store.CreateOrder(ctx, arg)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		stock := db.UpdateProductStockByIDParams{
			StockQuantity: oldStock - int32(req.Quantity),
			ID:            req.ProductID,
		}
		server.store.UpdateProductStockByID(ctx, stock)

		ctx.JSON(http.StatusOK, order)
		return

	}
	ctx.JSON(http.StatusConflict, gin.H{"error": "product is not enough"})

}

func TotalPrice(price_each, quantity int) int {
	return price_each * quantity
}
