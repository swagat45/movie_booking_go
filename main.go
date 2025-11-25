package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/files"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// "your/module/path/docs" // after running swag init

	_ "movie-booking-go/docs"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// @title Movie Ticket Booking API (Go)
// @version 1.0
// @description Movie booking backend intern assignment implemented in Go.
// @BasePath /api

func initDB() *gorm.DB {
	database, err := gorm.Open(sqlite.Open("movie_booking.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	if err := AutoMigrateModels(database); err != nil {
		log.Fatal("failed to run migrations: ", err)
	}

	// seed sample data (optional)
	seed(database)
	return database
}

func seed(database *gorm.DB) {
	// only seed if no movies
	var count int64
	database.Model(&Movie{}).Count(&count)
	if count > 0 {
		return
	}

	movie := Movie{
		Title:           "Inception",
		DurationMinutes: 148,
	}
	database.Create(&movie)

	show := Show{
		MovieID:    movie.ID,
		ScreenName: "Screen 1",
		DateTime:   time.Now().Add(2 * time.Hour),
		TotalSeats: 50,
	}
	database.Create(&show)
}

func main() {
	db = initDB()

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Swagger
	// docs.SwaggerInfo.BasePath = "/api"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.POST("/signup", SignupHandler)
		api.POST("/login", LoginHandler)

		api.GET("/movies", ListMoviesHandler)
		api.GET("/movies/:id/shows", ListShowsForMovieHandler)

		// Protected routes
		secured := api.Group("/")
		secured.Use(AuthMiddleware())
		{
			secured.POST("/shows/:id/book", BookSeatHandler)
			secured.POST("/bookings/:id/cancel", CancelBookingHandler)
			secured.GET("/my-bookings", MyBookingsHandler)
		}
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
