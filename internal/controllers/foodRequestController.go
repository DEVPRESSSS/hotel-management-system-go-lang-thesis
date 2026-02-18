package controllers

import (
	"HMS-GO/internal/models"
	"HMS-GO/internal/models/dto"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *Server) CreateFoodPayment(ctx *gin.Context) {

	// ─────────────────────────────────────────
	// ENVIRONMENT
	// ─────────────────────────────────────────
	secretKey := os.Getenv("PAYMONGO_SECRET_KEY")
	if secretKey == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Secret key is invalid"})
		return
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8085"
	}

	// ─────────────────────────────────────────
	// PARSE REQUEST
	// ─────────────────────────────────────────
	var req dto.FoodCheckoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if len(req.Items) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cart is empty"})
		return
	}

	// ─────────────────────────────────────────
	// BUILD LINE ITEMS
	// ─────────────────────────────────────────
	var lineItems []map[string]interface{}
	totalAmount := decimal.NewFromInt(0)

	for _, item := range req.Items {

		if item.Qty <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid quantity"})
			return
		}

		var food models.Food
		if err := s.Db.Where("food_id = ?", item.Id).First(&food).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "food not found"})
			return
		}

		qtyDecimal := decimal.NewFromInt(int64(item.Qty))
		itemTotal := qtyDecimal.Mul(food.Price)

		totalAmount = totalAmount.Add(itemTotal)

		amountInCents := itemTotal.Mul(decimal.NewFromInt(100)).IntPart()

		lineItems = append(lineItems, map[string]interface{}{
			"amount":   amountInCents,
			"currency": "PHP",
			"name":     food.Name,
			"quantity": item.Qty,
		})
	}

	// ─────────────────────────────────────────
	// CREATE PAYMONGO BODY
	// ─────────────────────────────────────────
	checkoutBody := map[string]interface{}{
		"data": map[string]interface{}{
			"attributes": map[string]interface{}{
				"line_items":           lineItems,
				"payment_method_types": []string{"gcash", "qrph"},
				"success_url":          fmt.Sprintf("%s/food/payment/success?session_id={CHECKOUT_SESSION_ID}", baseURL),
				"cancel_url":           fmt.Sprintf("%s/guest/food/services", baseURL),
				"description":          "Food order checkout",
				"metadata": map[string]interface{}{
					"total_amount": totalAmount.String(),
					"cart_items":   req.Items,
				},
			},
		},
	}

	body, err := json.Marshal(checkoutBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request body"})
		return
	}

	// ─────────────────────────────────────────
	// SEND TO PAYMONGO
	// ─────────────────────────────────────────
	request, err := http.NewRequest("POST",
		"https://api.paymongo.com/v1/checkout_sessions",
		bytes.NewBuffer(body),
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(secretKey, "")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to call PayMongo"})
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		ctx.JSON(resp.StatusCode, gin.H{"error": string(respBody)})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse PayMongo response"})
		return
	}

	// ─────────────────────────────────────────
	// EXTRACT RESPONSE
	// ─────────────────────────────────────────
	data := result["data"].(map[string]interface{})
	attributes := data["attributes"].(map[string]interface{})

	ctx.JSON(http.StatusOK, gin.H{
		"checkout_url": attributes["checkout_url"],
		"session_id":   data["id"],
	})
}

func (s *Server) FoodPaymentSuccess(ctx *gin.Context) {
	sessionID := ctx.Query("session_id")
	if sessionID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing session"})
		return
	}
	ctx.HTML(http.StatusOK, "food_payment_success.html", gin.H{"session_id": sessionID})
}

func (s *Server) InsertOrder(ctx *gin.Context) {

	var req dto.FoodCheckoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if len(req.Items) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cart is empty"})
		return
	}

	userID := s.GetUserId(ctx)
	sessionID := uuid.New().String() // generate a local session ID since no PayMongo session

	for _, item := range req.Items {
		if item.Qty <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid quantity"})
			return
		}

		var food models.Food
		if err := s.Db.Where("food_id = ?", item.Id).First(&food).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "food not found"})
			return
		}

		qty := decimal.NewFromInt(int64(item.Qty))
		totalPrice := food.Price.Mul(qty)

		order := models.Orders{
			OrderId:       uuid.New().String(),
			SessionId:     sessionID,
			UserId:        userID,
			BookId:        req.BookId,
			ProductName:   food.Name,
			Qty:           item.Qty,
			Price:         food.Price.InexactFloat64(),
			TotalPrice:    totalPrice.InexactFloat64(),
			PaymentStatus: "paid",
			CreatedAt:     time.Now(),
		}

		if err := s.Db.Create(&order).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert order"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Orders created successfully"})
}
