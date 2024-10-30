package main

import (
	"log"
	"products-api-with-jwt/config"
	"products-api-with-jwt/controllers"
	_ "products-api-with-jwt/docs" // Import docs for Swagger
	"products-api-with-jwt/middlewares"
	"products-api-with-jwt/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Setup database (SQLite)
	db, err := config.SetupDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected successfully")

	// Initialize DB for services
	authService := services.NewAuthService(db)
	bookService := services.NewBookService(db)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	bookController := controllers.NewBookController(bookService, authService)

	// Initialize router
	r := gin.Default()
	r.Use(middlewares.LoggingMiddleware())
	r.Use(middlewares.RateLimiterMiddleware())

	// Endpoint login (does not require JWT authentication)
	auth := r.Group("/auth")
	auth.POST("/login", authController.Login)
	auth.POST("/logout", authController.Logout)

	// Other endpoints require JWT authentication
	protected := r.Group("/")
	protected.Use(middlewares.JWTAuthMiddleware(authService))

	// Product endpoints
	book := protected.Group("/books")
	book.GET("/", bookController.GetBooks)             // Get all books
	book.GET("/:id", bookController.GetBookByID)       // Get book by ID
	book.POST("/", bookController.CreateBook)          // Add new book
	book.DELETE("/:id", bookController.DeleteBook)     // Delete book
	book.PUT("/:id", bookController.UpdateBook)        // Update book
	book.GET("/borrow/:id", bookController.BorrowBook) // Borrow book
	book.GET("/return/:id", bookController.ReturnBook) // Return book

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run server on port 8080
	r.Run(":8080")
}
