package handlers

import (
	"clothingretail/db"
	"clothingretail/models"
	"clothingretail/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RentClothing handles renting clothing items
func RentClothing(c *gin.Context) {
	var req models.RentalRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse dates
	dateBegin, err := parseDateTime(req.RentDateBegin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rent begin date format. Supported formats: dd-MMM-yyyy HH:mm, YYYY-MM-DD, ISO 8601"})
		return
	}

	dateEnd, err := parseDateTime(req.RentDateEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rent end date format. Supported formats: dd-MMM-yyyy HH:mm, YYYY-MM-DD, ISO 8601"})
		return
	}

	now := time.Now()
	rentalID := utils.GenerateID()

	// Insert rental record
	_, err = db.DB.Exec(
		`INSERT INTO clothing_rental (id, id_clothing_category_sub, id_clothing_size, id_clothing_customer, 
         clothes_qty_rent, clothes_qty_return, clothes_rent_date_begin, clothes_rent_date_end, 
         clothes_rent_date_actual_pickup, clothes_rent_date_actual_return, clothes_cat_status_sub, 
         created_at, updated_at) VALUES (?, ?, ?, ?, ?, 0, ?, ?, ?, '0001-01-01', 1, ?, ?)`,
		rentalID, req.IDClothingCategorySub, req.IDClothingSize, req.IDClothingCustomer, req.ClothesQtyRent,
		dateBegin, dateEnd, now, now, now,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update inventory movement
	inventoryID := utils.GenerateID()
	_, err = db.DB.Exec(
		`INSERT INTO clothing_inventory_movement (id, id_clothing_category, id_clothing_size, 
         clothes_movement_action, clothes_qty_in, clothes_qty_out, clothes_qty_total, 
         clothes_cat_status_sub, created_at, updated_at) 
         VALUES (?, ?, ?, 1, 0, ?, (SELECT COALESCE(SUM(clothes_qty_in - clothes_qty_out), 0) FROM clothing_inventory_movement 
         WHERE id_clothing_category = ? AND id_clothing_size = ?), 1, ?, ?)`,
		inventoryID, req.IDClothingCategorySub, req.IDClothingSize, req.ClothesQtyRent,
		req.IDClothingCategorySub, req.IDClothingSize, now, now,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Rental created successfully",
		"rental_id": rentalID,
	})
}

// ReturnClothing handles returning rented clothing items
func ReturnClothing(c *gin.Context) {
	var req models.ReturnRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get rental information
	var rental models.ClothingRental
	err := db.DB.QueryRow(
		`SELECT id, id_clothing_category_sub, id_clothing_size, id_clothing_customer, 
         clothes_qty_rent, clothes_qty_return FROM clothing_rental WHERE id = ?`,
		req.RentalID,
	).Scan(&rental.ID, &rental.IDClothingCategorySub, &rental.IDClothingSize,
		&rental.IDClothingCustomer, &rental.ClothesQtyRent, &rental.ClothesQtyReturn)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rental not found"})
		return
	}

	// Check if return quantity is valid
	remainingQty := rental.ClothesQtyRent - rental.ClothesQtyReturn
	if req.ClothesQtyReturn > remainingQty {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Return quantity exceeds rented quantity"})
		return
	}

	now := time.Now()
	newReturnQty := rental.ClothesQtyReturn + req.ClothesQtyReturn

	// Determine new status
	newStatus := 1 // active
	if newReturnQty == rental.ClothesQtyRent {
		newStatus = 3 // fully returned
	}

	// Update rental record
	_, err = db.DB.Exec(
		`UPDATE clothing_rental SET clothes_qty_return = ?, clothes_rent_date_actual_return = ?, 
         clothes_cat_status_sub = ?, updated_at = ? WHERE id = ?`,
		newReturnQty, now, newStatus, now, req.RentalID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update inventory movement
	inventoryID := utils.GenerateID()
	_, err = db.DB.Exec(
		`INSERT INTO clothing_inventory_movement (id, id_clothing_category, id_clothing_size, 
         clothes_movement_action, clothes_qty_in, clothes_qty_out, clothes_qty_total, 
         clothes_cat_status_sub, created_at, updated_at) 
         VALUES (?, ?, ?, 2, ?, 0, (SELECT COALESCE(SUM(clothes_qty_in - clothes_qty_out), 0) FROM clothing_inventory_movement 
         WHERE id_clothing_category = ? AND id_clothing_size = ?), 1, ?, ?)`,
		inventoryID, rental.IDClothingCategorySub, rental.IDClothingSize, req.ClothesQtyReturn,
		rental.IDClothingCategorySub, rental.IDClothingSize, now, now,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Return processed successfully",
		"status":  newStatus,
	})
}

// GetRentals retrieves all rentals
func GetRentals(c *gin.Context) {
	customerID := c.Query("customer_id")
	status := c.Query("status")

	query := `SELECT id, id_clothing_category_sub, id_clothing_size, id_clothing_customer, 
              clothes_qty_rent, clothes_qty_return, clothes_rent_date_begin, clothes_rent_date_end, 
              clothes_rent_date_actual_pickup, clothes_rent_date_actual_return, clothes_rent_status, 
              created_at, updated_at FROM clothing_rental WHERE 1=1`

	var args []interface{}

	if customerID != "" {
		query += " AND id_clothing_customer = ?"
		args = append(args, customerID)
	}

	if status != "" {
		query += " AND clothes_cat_status_sub = ?"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var rentals []models.ClothingRental
	for rows.Next() {
		var rental models.ClothingRental
		if err := rows.Scan(&rental.ID, &rental.IDClothingCategorySub, &rental.IDClothingSize,
			&rental.IDClothingCustomer, &rental.ClothesQtyRent, &rental.ClothesQtyReturn,
			&rental.ClothesRentDateBegin, &rental.ClothesRentDateEnd, &rental.ClothesRentDateActualPickup,
			&rental.ClothesRentDateActualReturn, &rental.ClothesRentStatus, &rental.CreatedAt,
			&rental.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		rentals = append(rentals, rental)
	}

	c.JSON(http.StatusOK, rentals)
}

// GetSizes retrieves all sizes, optionally filtered by subcategory
func GetSizes(c *gin.Context) {
	subcategoryID := c.Query("subcategory_id")

	query := `SELECT id, id_clothing_category_sub, clothes_size_name, clothes_size_notes, 
              clothes_size_status, created_at, updated_at 
              FROM clothing_size WHERE clothes_size_status = 1`

	var args []interface{}

	if subcategoryID != "" {
		query += " AND id_clothing_category_sub = ?"
		args = append(args, subcategoryID)
	}

	query += " ORDER BY clothes_size_name"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// Initialize with empty slice instead of nil
	sizes := []struct {
		ID                    int       `json:"id"`
		IDClothingCategorySub int       `json:"id_clothing_category_sub"`
		ClothesSizeName       string    `json:"clothes_size_name"`
		ClothesSizeNotes      string    `json:"clothes_size_notes"`
		ClothesSizeStatus     int       `json:"clothes_size_status"`
		CreatedAt             time.Time `json:"created_at"`
		UpdatedAt             time.Time `json:"updated_at"`
	}{}

	for rows.Next() {
		var size struct {
			ID                    int       `json:"id"`
			IDClothingCategorySub int       `json:"id_clothing_category_sub"`
			ClothesSizeName       string    `json:"clothes_size_name"`
			ClothesSizeNotes      string    `json:"clothes_size_notes"`
			ClothesSizeStatus     int       `json:"clothes_size_status"`
			CreatedAt             time.Time `json:"created_at"`
			UpdatedAt             time.Time `json:"updated_at"`
		}
		if err := rows.Scan(&size.ID, &size.IDClothingCategorySub, &size.ClothesSizeName,
			&size.ClothesSizeNotes, &size.ClothesSizeStatus, &size.CreatedAt,
			&size.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		sizes = append(sizes, size)
	}

	c.JSON(http.StatusOK, sizes)
}

// Helper function to parse multiple date/datetime formats
func parseDateTime(dateStr string) (time.Time, error) {
	// Try different formats
	formats := []string{
		"02-Jan-2006 15:04", // dd-MMM-yyyy HH:mm format (new format)
		"2006-01-02T15:04",  // HTML datetime-local format
		"2006-01-02 15:04",  // Common datetime format
		"2006-01-02",        // Date only
		time.RFC3339,        // ISO 8601
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
