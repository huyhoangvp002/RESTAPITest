package api

import (
	db "RESTAPITest/db/sqlc"
	"RESTAPITest/token"
	"RESTAPITest/util"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Address     string `json:"address" binding:"required"`
}

func (server *Server) Register(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	HashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateAccountParams{
		Username:     req.Username,
		HashPassword: HashedPassword,
		Role:         "buyer",
	}
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accountID, err := server.store.GetAccountIDByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg1 := db.CreateAccountInfoParams{
		AccountID:   accountID,
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
	}
	_, err = server.store.CreateAccountInfo(ctx, arg1)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)

}

type loginRequest struct {
	Username     string `json:"username" binding:"required"`
	HashPassword string `json:"password" binding:"required,min=6"`
}

type loginResponse struct {
	LoginToken string `json:"login_token"`
	UserName   string `json:"username" binding:"required"`
	Role       string `json:"role"`
}

func (server *Server) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account, err := server.store.GetAccountByUsername(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = util.CheckPassword(req.HashPassword, account.HashPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	loginToken, err := server.tokenMaker.CreateToken(
		req.Username,
		account.Role,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.SetCookie(
		"access_token", // tên cookie
		loginToken,     // giá trị: JWT
		int(server.config.AccessTokenDuration.Seconds()), // thời gian sống
		"/",   // path
		"",    //
		false, // secure
		false, //
	)

	rsp := loginResponse{
		LoginToken: loginToken,
		UserName:   req.Username,
		Role:       account.Role,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) GetAccountByUsername(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	account, err := server.store.GetAccountByUsername(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) ListAccounts(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Role == "admin" {

		account, err := server.store.ListAccounts(ctx, arg)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, account)
		return
	}
	account, err := server.store.GetAccountByUsername(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)

}

type updateRoleRequest struct {
	ID   int    `json:"id" binding:"required"`
	Role string `json:"role" binding:"required"`
}

func (server *Server) UpdateRole(ctx *gin.Context) {
	var req updateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !util.IsValidRole(req.Role) {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "Invalid Role"})
		return
	}

	arg := db.UpdateRoleParams{
		ID:   int64(req.ID),
		Role: req.Role,
	}

	account, err := server.store.UpdateRole(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)

}

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) DeleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload.Role == "admin" {

		_, err := server.store.GetAccountByID(ctx, req.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		server.store.DeleteAccount(ctx, req.ID)

		ctx.JSON(http.StatusOK, gin.H{"msg": "Delete Successfully"})
		return

	}

	accountID, err := server.store.GetIDByUserName(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if accountID != req.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": "Permission deny!"})
		return
	}
	server.store.DeleteAccount(ctx, req.ID)

	ctx.JSON(http.StatusOK, gin.H{"msg": "Delete Successfully"})

}

func (server *Server) LogOut(ctx *gin.Context) {
	ctx.SetCookie(
		"access_token", "", -1, "/", "", false, true,
	)
	ctx.JSON(http.StatusOK, gin.H{"msg": "see you again"})
}

func (server *Server) redirect(ctx *gin.Context) {
	ctx.Redirect(http.StatusTemporaryRedirect, "/static/landing_page.html")
}

type AccountByIDRequest struct {
	ID int `form:"id" binding:"required"`
}

func (server *Server) AccountByID(ctx *gin.Context) {
	var req AccountByIDRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fmt.Println("[DEBUG] ID:", req.ID)
	account, err := server.store.GetAccountByID(ctx, int64(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}
