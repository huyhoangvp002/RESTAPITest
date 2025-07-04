package api

import (
	db "RESTAPITest/db/sqlc"
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

	arg := db.CreateProductParams{
		Name:       req.Name,
		Price:      req.Price,
		CategoryID: sql.NullInt32{Int32: req.CategoryID, Valid: true},
		Value:      req.Value,
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
			ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

type getProductByCateRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getProductByCateRequest(ctx *gin.Context) {
	var req getProductByCateRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	product, err := server.store.ListProductsByCategoryID(ctx, sql.NullInt32{Int32: req.ID, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}
