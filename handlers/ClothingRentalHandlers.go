package handlers

import (
	"clothingretail/db"
	"clothingretail/models"
	"clothingretail/utils"
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
	dateBegin, err := time.Parse("2006-01-02", req.RentDateBegin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rent begin date format. Use YYYY-MM-DD"})
		return
	}

	dateEnd, err := time.Parse("2006-01-02", req.RentDateEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rent end date format. Use YYYY-MM-DD"})
		return
	}

	now := time.Now()
	rentalID := utils.GenerateID()

	// Insert rental record
	_, err = db.DB.Exec(
		`INSERT INTO clothing_rental (id, id_clothing_category_sub, id_clothing_size, id_clothing_customer, 
         clothes_qty_rent, clothes_qty_return, clothes_rent_date_begin, clothes_rent_date_end, 
         clothes_rent_date_actual_pickup, clothes_rent_date_actual_return, clothes_rent_status, 
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
         clothes_rent_status, created_at, updated_at) 
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
         clothes_rent_status = ?, updated_at = ? WHERE id = ?`,
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
         clothes_rent_status, created_at, updated_at) 
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
