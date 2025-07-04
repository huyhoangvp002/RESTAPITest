package api

import (
	db "RESTAPITest/db/sqlc"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createCategoryRespond struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (server *Server) createCategory(ctx *gin.Context) {
	var req createCategoryRespond
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("JSON bind err")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error"})
	}
	arg := db.CreateCategoryParams{
		Name: req.Name,
		Type: req.Type,
	}

	cate, err := server.store.CreateCategory(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error"})
	}
	ctx.JSON(http.StatusOK, cate)
}
