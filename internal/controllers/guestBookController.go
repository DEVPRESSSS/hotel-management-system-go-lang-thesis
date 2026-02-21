package controllers

import (
	"HMS-GO/internal/models"
	"HMS-GO/internal/models/dto"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v76"
	checkoutsession "github.com/stripe/stripe-go/v76/checkout/session"
	"gorm.io/gorm"
)

func (s *Server) FetchCalendar(ctx *gin.Context) {
	var books []models.Book

	roomId := ctx.Param("room_id")

	// Only fetch future/active bookings for calendar display
	if err := s.Db.
		Where("check_out_date >= ? AND room_id = ?", time.Now(), roomId).
		Order("check_in_date ASC").
		Find(&books).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch bookings",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"books": books,
	})
}

func (s *Server) RoomSelected(ctx *gin.Context) {
	roomId := ctx.Param("roomid")
	var room models.Room
	if err := s.Db.
		Preload("RoomType").
		Preload("Floor").
		Preload("Amenities").
		Where("room_id = ?", &roomId).
		First(&room).Error; err != nil {

		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Error fetching data!!!"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"room": room})
}

// Calculate boooking price (pre-booking)
func (s *Server) CalculateBookingPrice(ctx *gin.Context) {

	var req dto.PriceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Get room price
	var room models.Room
	if err := s.Db.Where("room_id = ?", req.RoomID).First(&room).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	layout := "2006-01-02 3:04 PM"

	checkIn, err := time.Parse(layout, req.CheckIn)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid check-in date"})
		return
	}

	checkOut, err := time.Parse(layout, req.CheckOut)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid check-out date"})
		return
	}

	if !checkOut.After(checkIn) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "checkout must be after checkin"})
		return
	}
	duration := checkOut.Sub(checkIn)
	nights := int(math.Ceil(duration.Hours() / 24))

	nightsDecimal := decimal.NewFromInt(int64(nights))
	numberGuest := decimal.NewFromInt(int64(req.Guest))
	total := nightsDecimal.Mul(room.Price).Mul(numberGuest)
	ctx.JSON(http.StatusOK, gin.H{
		"price_per_night": room.Price,
		"nights":          nights,
		"total":           total,
	})
}

// Submit booking
func (s *Server) ConfirmBooking(ctx *gin.Context) {

	userId, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var book models.Book

	if err := ctx.ShouldBindJSON(&book); err != nil {
		log.Print("This is the bad request", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate booking ID
	bookingID, err := GenerateBookingID(s.Db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate booking ID"})
		return
	}

	//Recalculate the total
	var room models.Room
	if err := s.Db.Where("room_id = ?", book.RoomId).First(&room).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	checkIn := book.CheckInDate
	checkOut := book.CheckOutDate

	if checkIn.IsZero() || checkOut.IsZero() {
		log.Print("This is the bad request")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "dates are required"})
		return
	}

	// Calculate actual price
	nights := int(checkOut.Sub(checkIn).Hours() / 24)
	nightsDecimal := decimal.NewFromInt(int64(nights))
	numberGuest := decimal.NewFromInt(int64(book.NumGuests))
	total := nightsDecimal.Mul(room.Price).Mul(numberGuest)

	// Assign server-controlled values
	book.BookId = bookingID
	book.UserId = userId.(string)
	book.PricePerNight = room.Price.InexactFloat64()
	book.TotalPrice = total.InexactFloat64()
	book.PaymentStatus = "Paid"

	// Assign BookId to each guest
	for i := range book.Guests {
		book.Guests[i].Id = fmt.Sprintf("BKGUEST-%03d", i+1)
		book.Guests[i].BookId = bookingID
		book.Guests[i].GuestNumber = i + 1
	}

	// Save booking
	if err := s.Db.Create(&book).Error; err != nil {
		log.Print("This is the bad request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":     "Booking confirmation has been sent to your email!",
		"room_price":  room.Price,
		"total_price": total,
	})
}

// Generate auto IncrementId
func GenerateBookingID(db *gorm.DB) (string, error) {
	var lastID string

	err := db.
		Model(&models.Book{}).
		Select("book_id").
		Order("book_id DESC").
		Limit(1).
		Scan(&lastID).Error

	if err != nil {
		return "", err
	}

	nextNumber := 1

	if lastID != "" {
		fmt.Sscanf(lastID, "BOOKING-%d", &nextNumber)
		nextNumber++
	}

	return fmt.Sprintf("BOOKING-%03d", nextNumber), nil
}

func GetPublishApiKey(ctx *gin.Context) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Env file failed to load")
	}

	ctx.JSON(http.StatusOK, gin.H{"publishableky": os.Getenv("STRIPE_PUBLISHABLE_KEY")})
}

func (s *Server) CreateCheckoutSession(ctx *gin.Context) {

	var req dto.PriceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	stripe.Key = os.Getenv("SECRET_STRIPE_KEY")

	// Get room details
	var room models.Room
	if err := s.Db.Preload("RoomType").Where("room_id = ?", req.RoomID).First(&room).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	layout := "2006-01-02 3:04 PM"

	checkIn, err := time.Parse(layout, req.CheckIn)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid check-in date"})
		return
	}

	checkOut, err := time.Parse(layout, req.CheckOut)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid check-out date"})
		return
	}

	if !checkOut.After(checkIn) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "checkout must be after checkin"})
		return
	}

	if !checkOut.After(checkIn) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "checkout must be after checkin"})
		return
	}

	duration := checkOut.Sub(checkIn)
	nights := int(math.Ceil(duration.Hours() / 24))
	nightsDecimal := decimal.NewFromInt(int64(nights))
	numberGuest := decimal.NewFromInt(int64(req.Guest))
	total := nightsDecimal.Mul(room.Price).Mul(numberGuest)

	// Convert to cents (Stripe uses smallest currency unit)
	amountInCents := total.Mul(decimal.NewFromInt(100)).IntPart()

	// Get base URL
	baseURL := os.Getenv("BASE_URL")
	// if baseURL == "" {
	// 	baseURL = "http://localhost:8085"
	// }

	// Get room type name (with fallback)
	roomTypeName := "Room"
	if room.RoomType.RoomTypeName != "" {
		roomTypeName = room.RoomType.RoomTypeName
	}

	// Create Checkout Session
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("php"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String(fmt.Sprintf("%s -%s", roomTypeName, req.RoomID)),
						Description: stripe.String(fmt.Sprintf("%d night • %d guest • %s to %s", nights, req.Guest, req.CheckIn, req.CheckOut)),
					},
					UnitAmount: stripe.Int64(amountInCents),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(fmt.Sprintf("%s/booking/success?session_id={CHECKOUT_SESSION_ID}", baseURL)),
		CancelURL:  stripe.String(fmt.Sprintf("%s/booking/summary", baseURL)),
		Metadata: map[string]string{
			"room_id":   req.RoomID,
			"check_in":  req.CheckIn,
			"check_out": req.CheckOut,
			"guests":    fmt.Sprintf("%d", req.Guest),
			"nights":    fmt.Sprintf("%d", nights),
		},
	}

	session, err := checkoutsession.New(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("Stripe session creation failed: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"url": session.URL,
	})
}

func (s *Server) BookingSuccess(ctx *gin.Context) {
	sessionID := ctx.Query("session_id")

	ctx.HTML(http.StatusOK, "booking_success.html", gin.H{
		"session_id": sessionID,
	})
}

// Create intent in paymongo
func (s *Server) CreatePaymentIntent(ctx *gin.Context) {

	// Get secret key
	secretKey := os.Getenv("PAYMONGO_SECRET_KEY")
	if secretKey == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Secret key is invalid"})
		return
	}
	//Get the base URL
	baseURL := os.Getenv("BASE_URL")
	// if baseURL == "" {
	// 	baseURL = "http://localhost:8085"
	// }
	// Parse request body to get booking details
	var req dto.PriceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// stripe.Key = os.Getenv("SECRET_STRIPE_KEY")

	// Get room details
	var room models.Room
	if err := s.Db.Preload("RoomType").Where("room_id = ?", req.RoomID).First(&room).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	layout := "2006-01-02 3:04 PM"

	checkIn, err := time.Parse(layout, req.CheckIn)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid check-in date"})
		return
	}

	checkOut, err := time.Parse(layout, req.CheckOut)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid check-out date"})
		return
	}

	if !checkOut.After(checkIn) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "checkout must be after checkin"})
		return
	}

	if !checkOut.After(checkIn) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "checkout must be after checkin"})
		return
	}

	duration := checkOut.Sub(checkIn)
	nights := int(math.Ceil(duration.Hours() / 24))
	nightsDecimal := decimal.NewFromInt(int64(nights))
	numberGuest := decimal.NewFromInt(int64(req.Guest))
	total := nightsDecimal.Mul(room.Price).Mul(numberGuest)

	// Convert to cents (Stripe uses smallest currency unit)
	amountInCents := total.Mul(decimal.NewFromInt(100)).IntPart()
	// Create checkout session
	checkoutURL := "https://api.paymongo.com/v1/checkout_sessions"

	checkoutBody := map[string]interface{}{
		"data": map[string]interface{}{
			"attributes": map[string]interface{}{
				"line_items": []map[string]interface{}{
					{
						"amount":   amountInCents,
						"currency": "PHP",
						"name":     "Booking/Reservation",
						"quantity": 1,
					},
				},
				//"payment_method_types": []string{"gcash", "paymaya", "card"},
				"payment_method_types": []string{"qrph", "gcash"},
				"success_url":          fmt.Sprintf("%s/booking/success?session_id={CHECKOUT_SESSION_ID}", baseURL),
				"cancel_url":           fmt.Sprintf("%s/booking/summary", baseURL),
				"description":          fmt.Sprintf("%s from %s to %s", req.RoomID, req.CheckIn, req.CheckOut),
				"metadata": map[string]interface{}{
					"room_id":   req.RoomID,
					"check_in":  req.CheckIn,
					"check_out": req.CheckOut,
					"guest":     req.Guest,
				},
			},
		},
	}

	// Marshal request body
	body, err := json.Marshal(checkoutBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create request body"})
		return
	}

	// Create HTTP request
	request, err := http.NewRequest("POST", checkoutURL, bytes.NewBuffer(body))
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create request: " + err.Error()})
		return
	}

	// Set headers and authentication
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(secretKey, "")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to call PayMongo: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	// Read response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to read response: " + err.Error()})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to parse response: " + err.Error()})
		return
	}

	// Check for errors
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		fmt.Printf("PayMongo Error Response: %s\n", string(bodyBytes))
		ctx.JSON(resp.StatusCode, result)
		return
	}

	// Extract checkout URL
	data := result["data"].(map[string]interface{})
	attributes := data["attributes"].(map[string]interface{})
	checkoutURLResponse := attributes["checkout_url"].(string)
	sessionID := data["id"].(string)

	ctx.JSON(http.StatusOK, gin.H{
		"checkout_url": checkoutURLResponse,
		"session_id":   sessionID,
	})
}
