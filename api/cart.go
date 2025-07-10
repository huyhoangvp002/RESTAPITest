package api

import (
	db "RESTAPITest/db/sqlc"
	"RESTAPITest/token"
	"database/sql"
	"log"
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

type ID_Request struct {
	ID int `uri:"id" binding:"required"`
}

type updateValueRequest struct {
	Value int `json:"value" binding:"required,min=1"`
}

// func (server *Server) UpdateProductInCart(ctx *gin.Context) {
// 	var req ID_Request
// 	if err := ctx.ShouldBindUri(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

// 	accountID_request, err := server.store.GetAccountIDByCartID(ctx, int64(req.ID))
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	accountID_payload, err := server.store.GetAccountIDByUsername(ctx, authPayload.Username)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	if accountID_request != int32(accountID_payload) {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"err": "Permission deny!"})
// 	}

// 	var req1 updateValueRequest
// 	if err := ctx.ShouldBindJSON(&req1); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	arg := db.UpdateCartValueParams{
// 		Value: int32(req1.Value),
// 	}

// 	cart, err := server.store.UpdateCartValue(ctx, arg)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, cart)
// }

func (server *Server) UpdateProductInCart(ctx *gin.Context) {
	log.Println("[DEBUG] UpdateProductInCart called")

	// 1. Bind URI param
	var req ID_Request
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Printf("[ERROR] Failed to bind URI: %v", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Printf("[DEBUG] URI param ID: %d", req.ID)

	// 2. Get Auth payload
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	log.Printf("[DEBUG] Auth username: %s", authPayload.Username)

	// 3. Get accountID by CartID
	accountID_request, err := server.store.GetAccountIDByCartID(ctx, int64(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[WARN] No cart found with id %d", req.ID)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		log.Printf("[ERROR] DB error GetAccountIDByCartID: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	log.Printf("[DEBUG] AccountID from cart: %d", accountID_request)

	// 4. Get accountID from auth username
	accountID_payload, err := server.store.GetAccountIDByUsername(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[WARN] No account found with username %s", authPayload.Username)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		log.Printf("[ERROR] DB error GetAccountIDByUsername: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	log.Printf("[DEBUG] AccountID from token: %d", accountID_payload)

	// 5. Check ownership
	if accountID_request != int32(accountID_payload) {
		log.Printf("[WARN] Permission denied! AccountID in cart: %d, from token: %d", accountID_request, accountID_payload)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Permission denied!"})
		return
	}

	// 6. Bind JSON body
	var req1 updateValueRequest
	if err := ctx.ShouldBindJSON(&req1); err != nil {
		log.Printf("[ERROR] Failed to bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Printf("[DEBUG] JSON body Value: %d", req1.Value)

	// 7. Build update param
	arg := db.UpdateCartValueParams{
		ID:    int64(req.ID),
		Value: int32(req1.Value),
	}
	log.Printf("[DEBUG] UpdateCartValueParams: %+v", arg)

	// 8. Call store
	cart, err := server.store.UpdateCartValue(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[WARN] No cart to update with id %d", req.ID)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		log.Printf("[ERROR] DB error UpdateCartValue: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	log.Printf("[DEBUG] Updated cart: %+v", cart)
	ctx.JSON(http.StatusOK, cart)
}

func (server *Server) DeleteCart(ctx *gin.Context) {
	var req ID_Request
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	accountID_request, err := server.store.GetAccountIDByCartID(ctx, int64(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accountID_payload, err := server.store.GetAccountIDByUsername(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if accountID_request != int32(accountID_payload) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": "Permission deny!"})
		return
	}

	err = server.store.DeleteCart(ctx, int64(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Delete Successfully!"})

}
