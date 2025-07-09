package api

import (
	db "RESTAPITest/db/sqlc"
	"fmt"

	"RESTAPITest/token"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createProductRequest struct {
	Name       string `json:"name" binding:"required"`
	Price      int32  `json:"price" binding:"required,min=1"`
	CategoryID int32  `json:"category_id" binding:"required"`
	Value      int32  `json:"value" binding:"required,min=0"`
}

func (server *Server) CreateProduct(ctx *gin.Context) {
	var req createProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//var errorr string

		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	id, err := server.store.GetIDByUserName(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	AccountID := sql.NullInt32{
		Int32: int32(id),
		Valid: true,
	}

	arg := db.CreateProductParams{
		Name:          req.Name,
		Price:         req.Price,
		DiscountPrice: req.Price,
		AccountID:     AccountID,
		CategoryID:    sql.NullInt32{Int32: req.CategoryID, Valid: true},
		Value:         req.Value,
	}
	product, err := server.store.CreateProduct(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

type getProductRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetProduct(ctx *gin.Context) {
	var req getProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	product, err := server.store.GetProduct(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

type getProductByCateRequest struct {
	CategoryName string `form:"name" binding:"required"`
	PageID       int32  `form:"page_id" binding:"required,min=1"`
	PageSize     int32  `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) GetProductByCate(ctx *gin.Context) {
	var req getProductByCateRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	cateID, err := server.store.GetCategoryIDByName(ctx, req.CategoryName)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	categoryID := sql.NullInt32{
		Int32: int32(cateID),
		Valid: true,
	}

	arg := db.ListProductsByCategoryIDParams{
		CategoryID: categoryID,
		Offset:     req.PageSize,
		Limit:      (req.PageID - 1) * req.PageSize,
	}

	products, err := server.store.ListProductsByCategoryID(ctx, arg)
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, products)

}

type updateProductRequest struct {
	ID    int64 `json:"id" binding:"required"`
	Price int32 `json:"price" binding:"required,min=1"`
	Value int32 `json:"value" binding:"required,min=0"`
}

func (server *Server) UdateProduct(ctx *gin.Context) {
	var req updateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	arg := db.UpdateProductParams{
		ID:    req.ID,
		Price: req.Price,
		Value: req.Value,
	}
	product, err := server.store.UpdateProduct(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, product)

}

type listByMaxPriceRequest struct {
	Price int32 `form:"price" binding:"required,min=1"`
}

func (server *Server) ListByMaxPrice(ctx *gin.Context) {
	var req listByMaxPriceRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	product, err := server.store.ListProductsByMaxPrice(ctx, req.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, product)

}

type listProductsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) ListProducts(ctx *gin.Context) {
	var req listProductsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListProductsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	product, err := server.store.ListProducts(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, product)
}

type listProductsByCustomerIDRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) ListProductByCustomerID(ctx *gin.Context) {
	var req listProductsByCustomerIDRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	cusID, err := server.store.GetAccountIDByUsername(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fmt.Println("[DEBUG]|Customer ID :", cusID)
	accountID := sql.NullInt32{
		Int32: int32(cusID),
		Valid: true,
	}

	arg := db.ListProductByAccountIDParams{
		AccountID: accountID,
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	}

	product, err := server.store.ListProductByAccountID(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, product)
}
