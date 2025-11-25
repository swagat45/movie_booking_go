package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB

// ========== AUTH DTOs ==========

type SignupRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// ========== BOOKING DTOs ==========

type BookSeatRequest struct {
	SeatNumber int `json:"seat_number" binding:"required"`
}

// ========== AUTH HANDLERS ==========

// POST /signup
func SignupHandler(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashed),
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use or invalid data"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// POST /login
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}

// ========== MOVIE & SHOW HANDLERS ==========

// GET /movies/
func ListMoviesHandler(c *gin.Context) {
	var movies []Movie
	if err := db.Find(&movies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}

	c.JSON(http.StatusOK, movies)
}

// GET /movies/:id/shows/
func ListShowsForMovieHandler(c *gin.Context) {
	movieID := c.Param("id")

	var shows []Show
	if err := db.Where("movie_id = ?", movieID).Preload("Movie").Find(&shows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shows"})
		return
	}

	c.JSON(http.StatusOK, shows)
}

// ========== BOOKING HANDLERS ==========

// POST /shows/:id/book/
func BookSeatHandler(c *gin.Context) {
	userID := GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	showID := c.Param("id")
	var req BookSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		var show Show
		if err := tx.First(&show, showID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return gin.Error{Err: errors.New("Show not found"), Type: gin.ErrorTypePublic}
			}
			return err
		}

		// Validate seat range
		if req.SeatNumber < 1 || req.SeatNumber > show.TotalSeats {
			return gin.Error{Err: errors.New("Invalid seat number"), Type: gin.ErrorTypePublic}
		}

		// Check overbooking
		var count int64
		if err := tx.Model(&Booking{}).
			Where("show_id = ? AND status = ?", show.ID, BookingStatusBooked).
			Count(&count).Error; err != nil {
			return err
		}
		if count >= int64(show.TotalSeats) {
			return gin.Error{Err: errors.New("Show is fully booked"), Type: gin.ErrorTypePublic}
		}

		// Check double booking for same seat
		var existing Booking
		if err := tx.Where("show_id = ? AND seat_number = ? AND status = ?",
			show.ID, req.SeatNumber, BookingStatusBooked).
			First(&existing).Error; err != nil {

			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		} else {
			return gin.Error{Err: errors.New("Seat already booked"), Type: gin.ErrorTypePublic}
		}

		booking := Booking{
			UserID:     userID,
			ShowID:     show.ID,
			SeatNumber: req.SeatNumber,
			Status:     BookingStatusBooked,
			CreatedAt:  time.Now(),
		}

		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		// Attach to context via error for response after tx
		c.Set("new_booking", booking)
		return nil
	})

	if err != nil {
		if ginErr, ok := err.(gin.Error); ok && ginErr.Type == gin.ErrorTypePublic {
			c.JSON(http.StatusBadRequest, gin.H{"error": ginErr.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to book seat"})
		return
	}

	v, _ := c.Get("new_booking")
	booking := v.(Booking)
	c.JSON(http.StatusCreated, booking)
}

// POST /bookings/:id/cancel/
func CancelBookingHandler(c *gin.Context) {
	userID := GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	bookingID := c.Param("id")
	var booking Booking
	if err := db.First(&booking, bookingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch booking"})
		return
	}

	// Security: user cannot cancel others' bookings
	if booking.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to cancel this booking"})
		return
	}

	if booking.Status != BookingStatusBooked {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking is not in booked status"})
		return
	}

	if err := db.Model(&booking).Update("status", BookingStatusCancelled).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking cancelled successfully"})
}

// GET /my-bookings/
func MyBookingsHandler(c *gin.Context) {
	userID := GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var bookings []Booking
	if err := db.Where("user_id = ?", userID).
		Preload("Show").
		Preload("Show.Movie").
		Order("created_at DESC").
		Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}
