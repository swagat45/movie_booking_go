package main

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `json:"name"`
	Email        string         `gorm:"uniqueIndex" json:"email"`
	PasswordHash string         `json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Movie struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Title           string         `json:"title"`
	DurationMinutes int            `json:"duration_minutes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Shows           []Show         `json:"shows,omitempty"`
}

type Show struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	MovieID    uint           `json:"movie_id"`
	Movie      Movie          `json:"movie,omitempty"`
	ScreenName string         `json:"screen_name"`
	DateTime   time.Time      `json:"date_time"`
	TotalSeats int            `json:"total_seats"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Bookings []Booking `json:"bookings,omitempty"`
}

type BookingStatus string

const (
	BookingStatusBooked    BookingStatus = "booked"
	BookingStatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `json:"user_id"`
	User       User           `json:"user,omitempty"`
	ShowID     uint           `json:"show_id"`
	Show       Show           `json:"show,omitempty"`
	SeatNumber int            `json:"seat_number"`
	Status     BookingStatus  `gorm:"type:varchar(20);index" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// AutoMigrateModels migrates all models
func AutoMigrateModels(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Movie{}, &Show{}, &Booking{})
}
