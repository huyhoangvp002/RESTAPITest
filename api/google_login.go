package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	db "RESTAPITest/db/sqlc"
	"RESTAPITest/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func (server *Server) HandleGoogleLogin(ctx *gin.Context) {
	state := "randomStateString123"

	config := util.GetGoogleOAuthConfig()
	fmt.Printf("=== OAUTH DEBUG ===\n")
	fmt.Printf("ClientID: '%s'\n", config.ClientID)
	fmt.Printf("RedirectURL: '%s'\n", config.RedirectURL)

	url := config.AuthCodeURL(state)
	fmt.Printf("Generated URL: %s\n", url)
	fmt.Printf("==================\n")

	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (server *Server) HandleGoogleCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	token, err := util.GetGoogleOAuthConfig().Exchange(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
		return
	}

	// Tạo client và gọi API để lấy thông tin user
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	// Parse user info
	data, _ := io.ReadAll(resp.Body)
	var userInfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.Unmarshal(data, &userInfo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unmarshal failed"})
		return
	}

	ok, err := server.store.CheckEmailExists(ctx, userInfo.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if !ok {
		password := util.RandomString(8)
		hash_password, err := util.HashPassword(password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		arg := db.CreateAccountParams{
			Username:     util.RandomString(8),
			HashPassword: hash_password,
			Role:         "buyer",
		}
		account, err := server.store.CreateAccount(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		arg1 := db.CreateAccountInfoParams{
			Name:  userInfo.Name,
			Email: userInfo.Email,
		}

		_, err = server.store.CreateAccountInfo(ctx, arg1)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// Tạo JWT và gửi lại bằng cookies
		tokenString, err := server.tokenMaker.CreateToken(account.Username, account.Role, server.config.AccessTokenDuration)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "JWT error"})
			return
		}

		ctx.SetCookie("access_token", tokenString, 3600, "/", "", false, false)
		ctx.Redirect(http.StatusTemporaryRedirect, "/static/info.html")
		return
	}

	account_id, err := server.store.GetAccountIDByEmail(ctx, userInfo.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	account, err := server.store.GetAccountByID(ctx, account_id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	tokenString, err := server.tokenMaker.CreateToken(account.Username, account.Role, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "JWT error"})
		return
	}

	ctx.SetCookie("access_token", tokenString, 3600, "/", "", false, false)
	ctx.Redirect(http.StatusTemporaryRedirect, "/static/home.html")

}
