package main

import (
	"clothingretail/conf"
	"clothingretail/db"
	"clothingretail/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	conf.CheckRunMode()
	if err := db.InitDB("./db/clothing.db"); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.CloseDB()

	// Initialize Gin router
	router := gin.Default()

	// Serve static files (HTML, CSS, JS)
	router.StaticFile("/", "./templates/index.html")
	router.Static("/static", "./templates")

	// API routes
	api := router.Group("/api")
	{
		// Category routes
		api.POST("/categories", handlers.CreateCategory)
		api.GET("/categories", handlers.GetCategories)
		api.GET("/categories/:id", handlers.GetCategoryByID)
		api.PUT("/categories/:id", handlers.UpdateCategory)
		api.DELETE("/categories/:id", handlers.DeleteCategory)

		// Subcategory routes
		api.POST("/categories-sub", handlers.CreateCategorySub)
		api.GET("/categories-sub", handlers.GetCategoriesSub)
		api.GET("/categories-sub/:id", handlers.GetCategorySubByID)
		api.PUT("/categories-sub/:id", handlers.UpdateCategorySub)
		api.DELETE("/categories-sub/:id", handlers.DeleteCategorySub)

		// Customer routes
		api.POST("/customers", handlers.CreateCustomer)
		api.GET("/customers", handlers.GetCustomers)
		api.GET("/customers/:id", handlers.GetCustomerByID)

		// Rental routes
		api.POST("/rentals", handlers.RentClothing)
		api.POST("/rentals/return", handlers.ReturnClothing)
		api.GET("/rentals", handlers.GetRentals)
	}

	// Start server
	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
