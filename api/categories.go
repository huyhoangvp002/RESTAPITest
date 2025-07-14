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
		return
	}

	arg := db.CreateCategoryParams{
		Name:      req.Name,
		Type:      req.Type,
		AccountID: account_IDRaw,
	}

	cate, err := server.store.CreateCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error"})
		return
	}
	ctx.JSON(http.StatusOK, cate)
}

type listCategoriesRequest struct {
	PageID   int `form:"page_id" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) ListCategories(ctx *gin.Context) {
	var req listCategoriesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.ListCategoriesParams{
		Limit:  int32(req.PageID),
		Offset: int32((req.PageID - 1) * req.PageSize),
	}

	cate, err := server.store.ListCategories(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, cate)
}

func (server *Server) DeleteCategories(ctx *gin.Context) {
	var req IDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := server.store.DeleteCategory(ctx, req.ID)
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
