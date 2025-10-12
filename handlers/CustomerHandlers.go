package handlers

import (
	"clothingretail/db"
	"clothingretail/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateCustomer handles creating a new customer
func CreateCustomer(c *gin.Context) {
	var customer models.ClothingCustomer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()
	customer.CustStatus = 1

	result, err := db.DB.Exec(
		`INSERT INTO clothing_customer (cust_name, cust_address, cust_city, cust_phone, cust_email, 
         cust_notes, cust_status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		customer.CustName, customer.CustAddress, customer.CustCity, customer.CustPhone,
		customer.CustEmail, customer.CustNotes, customer.CustStatus, customer.CreatedAt, customer.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	customer.ID = int(id)

	c.JSON(http.StatusCreated, customer)
}

// GetCustomers retrieves all customers
func GetCustomers(c *gin.Context) {
	rows, err := db.DB.Query(
		`SELECT id, cust_name, cust_address, cust_city, cust_phone, cust_email, cust_notes, 
         cust_status, created_at, updated_at FROM clothing_customer WHERE cust_status = 1`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var customers []models.ClothingCustomer
	for rows.Next() {
		var customer models.ClothingCustomer
		if err := rows.Scan(&customer.ID, &customer.CustName, &customer.CustAddress, &customer.CustCity,
			&customer.CustPhone, &customer.CustEmail, &customer.CustNotes, &customer.CustStatus,
			&customer.CreatedAt, &customer.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		customers = append(customers, customer)
	}

	c.JSON(http.StatusOK, customers)
}

// GetCustomerByID retrieves a single customer by ID
func GetCustomerByID(c *gin.Context) {
	id := c.Param("id")
	var customer models.ClothingCustomer

	err := db.DB.QueryRow(
		`SELECT id, cust_name, cust_address, cust_city, cust_phone, cust_email, cust_notes, 
         cust_status, created_at, updated_at FROM clothing_customer WHERE id = ?`,
		id,
	).Scan(&customer.ID, &customer.CustName, &customer.CustAddress, &customer.CustCity,
		&customer.CustPhone, &customer.CustEmail, &customer.CustNotes, &customer.CustStatus,
		&customer.CreatedAt, &customer.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}
