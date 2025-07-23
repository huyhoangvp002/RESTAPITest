package api

import (
	db "RESTAPITest/db/sqlc"
	"RESTAPITest/token"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createCAccountInfoRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Address     string `json:"address" binding:"required"`
}

func (server *Server) CreateAccountInfo(ctx *gin.Context) {
	var req createCAccountInfoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
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

	arg := db.CreateAccountInfoParams{
		Name:        req.Name,
		AccountID:   id,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
	}

	customer, err := server.store.CreateAccountInfo(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

type updateInfoRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Address     string `json:"address" bindind:"required"`
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"omitempty,numeric,len=10"`
}

func (server *Server) UpdateAccountInfo(ctx *gin.Context) {
	var req updateInfoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	AuthPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	accountID, err := server.store.GetAccountIDByUsername(ctx, AuthPayload.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.UpdateAccountInfoByAccountIDParams{
		AccountID:   accountID,
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
	}
	account_info, err := server.store.UpdateAccountInfoByAccountID(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account_info)

}

func (server *Server) GetAccountInfo(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	accountID, err := server.store.GetAccountIDByUsername(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	account_info, err := server.store.GetAccountInfoByAccountID(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account_info)
}
