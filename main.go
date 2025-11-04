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
	if err := db.InitDB(conf.Koan.String(conf.RunMode + ".ds_sqlite")); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.CloseDB()

	// Create default user if none exists
	if err := handlers.CreateDefaultUser(); err != nil {
		log.Println("Warning: Failed to create default user:", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Serve static files (CSS, JS)
	router.Static("/static", "./templates")

	// Public routes (no authentication required)
	router.GET("/login", func(c *gin.Context) {
		c.File("./templates/login.html")
	})

	// Authentication API routes (public for login, but logout should be protected)
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", handlers.Login)
	}

	// Protected routes - apply middleware
	protected := router.Group("/")
	protected.Use(handlers.AuthMiddleware())
	{
		// Serve index page (protected)
		protected.GET("/", func(c *gin.Context) {
			log.Println("Come to index")
			c.File("./templates/index.html")
		})

		// HTML page routes (protected)
		protected.GET("/create-category", func(c *gin.Context) {
			log.Println("Creating category")
			c.File("./templates/create-category.html")
		})
		protected.GET("/edit-category", func(c *gin.Context) {
			c.File("./templates/edit-category.html")
		})
		protected.GET("/create-category-sub", func(c *gin.Context) {
			c.File("./templates/create-category-sub.html")
		})
		protected.GET("/edit-category-sub", func(c *gin.Context) {
			c.File("./templates/edit-category-sub.html")
		})
		protected.GET("/create-customer", func(c *gin.Context) {
			c.File("./templates/create-customer.html")
		})
		protected.GET("/create-rental", func(c *gin.Context) {
			c.File("./templates/create-rental.html")
		})
		protected.GET("/return-rental", func(c *gin.Context) {
			c.File("./templates/return-rental.html")
		})

		// Logout route (protected - must be authenticated to logout)
		protected.POST("/api/auth/logout", handlers.Logout)

		// API routes (protected)
		api := protected.Group("/api")
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

			// Size routes
			api.GET("/sizes", handlers.GetSizes)

			// Rental routes
			api.POST("/rentals", handlers.RentClothing)
			api.POST("/rentals/return", handlers.ReturnClothing)
			api.GET("/rentals", handlers.GetRentals)
		}
	}

	// Start server
	log.Printf("Server starting on %s\n", conf.Koan.String(conf.RunMode+".port_api"))
	if err := router.Run(conf.Koan.String(conf.RunMode + ".port_api")); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
