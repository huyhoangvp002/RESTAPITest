package api

import (
	db "RESTAPITest/db/sqlc"
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createShipmentRequest struct {
	OrderID int `json:"order_id" binding:"required"`
}

type addressRequest struct {
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	Address string `json:"address" binding:"required"`
}
type shipmentDeliveryRequest struct {
	FromAddress addressRequest `json:"from_address" binding:"required"`
	ToAddress   addressRequest `json:"to_address" binding:"required"`
	Fee         int            `json:"fee" binding:"required,min=0"`
}

type NullableString struct {
	String string `json:"String"`
	Valid  bool   `json:"Valid"`
}

type shipmentDeliveryResponse struct {
	ShipmentCode NullableString `json:"shipment_code"`
	Fee          int            `json:"fee"`
	Status       NullableString `json:"status"`
}

func (server *Server) CreateShipment(ctx *gin.Context) {
	var req createShipmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERROR] BindJSON error: %v", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Printf("[INFO] CreateShipment request: %+v", req)

	// Get seller ID
	from_id, err := server.store.GetSellerIDByOrderID(ctx, int64(req.OrderID))
	if err != nil {
		log.Printf("[ERROR] GetSellerIDByOrderID: %v", err)
		status := http.StatusInternalServerError
		if err == sql.ErrNoRows {
			status = http.StatusNotFound
		}
		ctx.JSON(status, errorResponse(err))
		return
	}
	log.Printf("[INFO] SellerID: %d", from_id)

	// Get buyer ID
	to_id, err := server.store.GetBuyerIDByOrderID(ctx, int64(req.OrderID))
	if err != nil {
		log.Printf("[ERROR] GetBuyerIDByOrderID: %v", err)
		status := http.StatusInternalServerError
		if err == sql.ErrNoRows {
			status = http.StatusNotFound
		}
		ctx.JSON(status, errorResponse(err))
		return
	}
	log.Printf("[INFO] BuyerID: %d", to_id)

	// Fetch seller details
	seller_name, err := server.store.GetNameForShipment(ctx, from_id)
	if err != nil {
		log.Printf("[ERROR] GetNameForShipment seller: %v", err)
		ctx.JSON(statusCodeForError(err), errorResponse(err))
		return
	}
	seller_phone, err := server.store.GetPhoneForShipment(ctx, from_id)
	if err != nil {
		log.Printf("[ERROR] GetPhoneForShipment seller: %v", err)
		ctx.JSON(statusCodeForError(err), errorResponse(err))
		return
	}
	seller_address, err := server.store.GetAddressForShipment(ctx, from_id)
	if err != nil {
		log.Printf("[ERROR] GetAddressForShipment seller: %v", err)
		ctx.JSON(statusCodeForError(err), errorResponse(err))
		return
	}

	from_address := addressRequest{
		Name:    seller_name,
		Phone:   seller_phone,
		Address: seller_address,
	}
	log.Printf("[INFO] FromAddress: %+v", from_address)

	// Fetch buyer details
	buyer_name, err := server.store.GetNameForShipment(ctx, to_id)
	if err != nil {
		log.Printf("[ERROR] GetNameForShipment buyer: %v", err)
		ctx.JSON(statusCodeForError(err), errorResponse(err))
		return
	}
	buyer_phone, err := server.store.GetPhoneForShipment(ctx, to_id)
	if err != nil {
		log.Printf("[ERROR] GetPhoneForShipment buyer: %v", err)
		ctx.JSON(statusCodeForError(err), errorResponse(err))
		return
	}
	buyer_address, err := server.store.GetAddressForShipment(ctx, to_id)
	if err != nil {
		log.Printf("[ERROR] GetAddressForShipment buyer: %v", err)
		ctx.JSON(statusCodeForError(err), errorResponse(err))
		return
	}

	to_address := addressRequest{
		Name:    buyer_name,
		Phone:   buyer_phone,
		Address: buyer_address,
	}
	log.Printf("[INFO] ToAddress: %+v", to_address)

	// Get fee
	fee, err := server.store.GetTotalPriceByID(ctx, int64(req.OrderID))
	if err != nil {
		log.Printf("[ERROR] GetTotalPriceByID: %v", err)
		ctx.JSON(statusCodeForError(err), errorResponse(err))
		return
	}
	log.Printf("[INFO] Fee: %d", fee)

	// Prepare payload
	shipment_delivery := shipmentDeliveryRequest{
		FromAddress: from_address,
		ToAddress:   to_address,
		Fee:         int(fee),
	}
	log.Printf("[DEBUG] Shipment payload: %+v", shipment_delivery)

	// Call delivery API
	jsonData, err := json.Marshal(shipment_delivery)
	if err != nil {
		log.Printf("[ERROR] JSON Marshal: %v", err)
		ctx.JSON(statusCodeForError(err), errorResponse(err))
		return
	}

	deliveryAPI_URL := "http://localhost:9999/api/shipment"
	apiKey := "1183b817ec21e4a8bc2c409cc17135c0"

	httpReq, err := http.NewRequest("POST", deliveryAPI_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[ERROR] NewRequest: %v", err)
		ctx.JSON(statusCodeForError(err), errorResponse(err))
		return
	}
	httpReq.Header.Set("Authorization", "ApiKey "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("[ERROR] HTTP client.Do: %v", err)
		ctx.JSON(statusCodeForError(err), errorResponse(err))
		return
	}
	defer resp.Body.Close()
	log.Printf("[INFO] Delivery API response status: %s", resp.Status)

	var deliveryResp shipmentDeliveryResponse
	if err := json.NewDecoder(resp.Body).Decode(&deliveryResp); err != nil {
		log.Printf("[ERROR] Failed to decode delivery response: %v", err)
		ctx.JSON(statusCodeForError(err), gin.H{"error": "Invalid response"})
		return
	}

	if deliveryResp.ShipmentCode.String == "" {
		log.Printf("[ERROR] shipment_code is missing in delivery response")
		ctx.JSON(statusCodeForError(err), gin.H{"error": "Missing shipment_code"})
		return
	}

	log.Printf("[INFO] Delivery created: %+v", deliveryResp)

	//load into DB

	arg := db.CreateShipmentParams{
		OrderID:      int64(req.OrderID),
		ShipmentCode: deliveryResp.ShipmentCode.String,
		Fee:          int64(deliveryResp.Fee),
		Status:       "created",
	}

	shipment, err := server.store.CreateShipment(ctx, arg)
	if err != nil {
		log.Printf("[ERROR] Save shipment to DB: %v", err)
		ctx.JSON(statusCodeForError(err), errorResponse(err))
		return
	}

	log.Printf("[INFO] Saved shipment in local DB: %+v", shipment)
	ctx.JSON(http.StatusOK, shipment)
}
