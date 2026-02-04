package controllers

import (
	"HMS-GO/internal/models"
	"HMS-GO/internal/models/dto"
	"fmt"
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

	// Parse dates
	// checkIn, err := time.Parse("2006-01-02", req.CheckIn)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid check-in date"})
	// 	return
	// }

	// checkOut, err := time.Parse("2006-01-02", req.CheckOut)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid check-out date"})
	// 	return
	// }
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

	// Load Stripe key
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
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

	// Calculate total (same logic as CalculateBookingPrice)

	duration := checkOut.Sub(checkIn)
	nights := int(math.Ceil(duration.Hours() / 24))
	nightsDecimal := decimal.NewFromInt(int64(nights))
	numberGuest := decimal.NewFromInt(int64(req.Guest))
	total := nightsDecimal.Mul(room.Price).Mul(numberGuest)

	// Convert to cents (Stripe uses smallest currency unit)
	amountInCents := total.Mul(decimal.NewFromInt(100)).IntPart()

	// Get base URL
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8085"
	}

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
