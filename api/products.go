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
	Stock      int32  `json:"stock" binding:"required,min=0"`
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

	arg := db.CreateProductParams{
		Name:          req.Name,
		Price:         req.Price,
		DiscountPrice: req.Price,
		AccountID:     id,
		CategoryID:    int64(req.CategoryID),
		StockQuantity: req.Stock,
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

	arg := db.ListProductsByCategoryIDParams{
		CategoryID: cateID,
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

type IDRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type updateProductRequest struct {
	Price int32 `json:"price" binding:"required,min=1"`
	Stock int32 `json:"stock" binding:"required,min=0"`
}

func (server *Server) UdateProduct(ctx *gin.Context) {
	var idReq IDRequest
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	authAccountID, err := server.store.GetAccountIDByUsername(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	reqAccountID, err := server.store.GetAccountIDbyProductID(ctx, idReq.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if reqAccountID != authAccountID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": "Permission deny!"})
	}

	var req updateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	arg := db.UpdateProductParams{
		ID:            idReq.ID,
		Price:         req.Price,
		StockQuantity: req.Stock,
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

	arg := db.ListProductByAccountIDParams{
		AccountID: cusID,
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

func (server *Server) DeleteProduct(ctx *gin.Context) {
	var req IDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role == "admin" {
		err := server.store.DeleteProduct(ctx, req.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": "Delete Successfully!"})
		return
	}

	authAccountID, err := server.store.GetAccountIDByUsername(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	reqAccountID, err := server.store.GetAccountIDbyProductID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if reqAccountID != authAccountID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": "Permission deny!"})
		return
	}

	err = server.store.DeleteProduct(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"err": "Delete Successfully!"})

}

type searchProductByNameRequest struct {
	Name string `form:"name"`
}

func (server *Server) SearchProductByName(ctx *gin.Context) {
	var pagingReq listProductsRequest
	if err := ctx.ShouldBindQuery(&pagingReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListProductsParams{
		Limit:  pagingReq.PageSize,
		Offset: (pagingReq.PageID - 1) * pagingReq.PageSize,
	}
	var nameReq searchProductByNameRequest
	ctx.ShouldBindQuery(&nameReq)
	if nameReq.Name == "" {
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
		return
	}

	productName := sql.NullString{
		String: nameReq.Name,
		Valid:  true,
	}

	arg2 := db.SearchProductsByNameParams{
		Column1: productName,
		Limit:   pagingReq.PageSize,
		Offset:  (pagingReq.PageID - 1) * pagingReq.PageSize,
	}
	product, err := server.store.SearchProductsByName(ctx, arg2)
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
