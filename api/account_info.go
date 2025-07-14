package api

import (
	db "RESTAPITest/db/sqlc"
	"RESTAPITest/token"
	"database/sql"
	"log"
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
	Email       *string `json:"email" binding:"omitempty,email"`
	Address     *string `json:"address"`
	Name        *string `json:"name"`
	PhoneNumber *string `json:"phone_number" binding:"omitempty,numeric,len=10"`
}

func (server *Server) UpdateAccountInfo(ctx *gin.Context) {
	log.Println("[DEBUG] UpdateAccountInfo: start handler")

	// 1. Get URI param
	var req IDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Printf("[ERROR] BindUri failed: %v\n", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Printf("[DEBUG] Parsed URI param: ID=%d\n", req.ID)

	// 2. Get auth payload
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	log.Printf("[DEBUG] Authenticated Username: %s\n", authPayload.Username)

	// 3. Get accountID of requester
	authID, err := server.store.GetAccountIDByUsername(ctx, authPayload.Username)
	if err != nil {
		log.Printf("[ERROR] GetAccountIDByUsername failed: %v\n", err)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	log.Printf("[DEBUG] Authenticated User ID: %d\n", authID)

	// 5. Get target accountID by request param
	reqAccountID, err := server.store.GetAccountID(ctx, req.ID)
	if err != nil {
		log.Printf("[ERROR] GetAccountID failed: %v\n", err)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	log.Printf("[DEBUG] Target Account ID from URI: %d\n", reqAccountID)

	// 6. Check permission
	if authID != reqAccountID {
		log.Printf("[WARN] Permission denied: requesterID=%d, targetID=%d\n", authID, req.ID)
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": "Permission Deny!"})
		return
	}
	log.Println("[DEBUG] Permission check passed")

	// 7. Parse JSON body
	var infoReq updateInfoRequest
	if err := ctx.ShouldBindJSON(&infoReq); err != nil {
		log.Printf("[ERROR] BindJSON failed: %v\n", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Printf("[DEBUG] Parsed JSON: %+v\n", infoReq)

	// 8. Update name if provided
	if infoReq.Name != nil {
		log.Printf("[DEBUG] Updating Name to: %s\n", *infoReq.Name)

		info := db.UpdateAccountInfoNameParams{
			ID:   authID,
			Name: *infoReq.Name,
		}

		if err := server.store.UpdateAccountInfoName(ctx, info); err != nil {
			log.Printf("[ERROR] DB update failed: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		log.Println("[INFO] Name updated successfully")
	}

	if infoReq.Email != nil {
		log.Printf("[DEBUG] Updating Email to: %s\n", *infoReq.Email)

		info := db.UpdateAccountInfoEmailParams{
			ID:    authID,
			Email: *infoReq.Email,
		}

		if err := server.store.UpdateAccountInfoEmail(ctx, info); err != nil {
			log.Printf("[ERROR] DB update failed: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		log.Println("[INFO] Email updated successfully")
	}

	if infoReq.PhoneNumber != nil {
		log.Printf("[DEBUG] Updating Phone Number to: %s\n", *infoReq.PhoneNumber)

		info := db.UpdateAccountInfoPhoneNumberParams{
			ID:          authID,
			PhoneNumber: *infoReq.PhoneNumber,
		}

		if err := server.store.UpdateAccountInfoPhoneNumber(ctx, info); err != nil {
			log.Printf("[ERROR] DB update failed: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		log.Println("[INFO] Phone Number updated successfully")
	}

	if infoReq.Address != nil {
		log.Printf("[DEBUG] Updating Address to: %s\n", *infoReq.Address)

		info := db.UpdateAccountInfoAddressParams{
			ID:      authID,
			Address: *infoReq.Address,
		}

		if err := server.store.UpdateAccountInfoAddress(ctx, info); err != nil {
			log.Printf("[ERROR] DB update failed: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		log.Println("[INFO] Address updated successfully")
	}

	// 9. Success response
	ctx.JSON(http.StatusOK, gin.H{"status": "updated"})
	log.Println("[DEBUG] UpdateAccountInfo: completed successfully")
}

func (server *Server) GetAccountInfo(ctx *gin.Context) {
	var req IDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload.Role == "admin" {
		account_info, err := server.store.GetAccountInfo(ctx, req.ID)
		if err != nil {
			if err != sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, account_info)
		return
	}

	accountID, err := server.store.GetAccountID(ctx, req.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if accountID != req.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": "Permission Deny!"})
		return
	}
	account_info, err := server.store.GetAccountInfo(ctx, req.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account_info)
}
