package handlers

import (
	"clothingretail/db"
	"clothingretail/models"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateCategorySub handles creating a new clothing subcategory
func CreateCategorySub(c *gin.Context) {
	var categorySub models.ClothingCategorySub

	if err := c.ShouldBindJSON(&categorySub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	categorySub.CreatedAt = time.Now()
	categorySub.UpdatedAt = time.Now()
	categorySub.ClothesCatStatusSub = 1

	result, err := db.DB.Exec(
		`INSERT INTO clothing_category_sub (id_clothing_category, clothes_cat_name_sub, clothes_cat_location_sub, 
         clothes_picture_1, clothes_picture_2, clothes_picture_3, clothes_picture_4, clothes_picture_5, 
         clothes_cat_status_sub, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		categorySub.IDClothingCategory, categorySub.ClothesCatNameSub, categorySub.ClothesCatLocationSub,
		categorySub.ClothesPicture1, categorySub.ClothesPicture2, categorySub.ClothesPicture3,
		categorySub.ClothesPicture4, categorySub.ClothesPicture5, categorySub.ClothesCatStatusSub,
		categorySub.CreatedAt, categorySub.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("Error inserting category sub: %v\n", err)
		return
	}

	id, _ := result.LastInsertId()
	categorySub.ID = int(id)

	c.JSON(http.StatusCreated, categorySub)
}

// GetCategoriesSub retrieves all clothing subcategories
func GetCategoriesSub(c *gin.Context) {
	categoryID := c.Query("category_id")

	query := `SELECT id, id_clothing_category, clothes_cat_name_sub, clothes_cat_location_sub, 
              clothes_picture_1, clothes_picture_2, clothes_picture_3, clothes_picture_4, clothes_picture_5, 
              clothes_cat_status_sub, created_at, updated_at FROM clothing_category_sub WHERE clothes_cat_status_sub = 1`

	var rows *sql.Rows
	var err error

	if categoryID != "" {
		query += " AND id_clothing_category = ?"
		rows, err = db.DB.Query(query, categoryID)
	} else {
		rows, err = db.DB.Query(query)
	}
	fmt.Printf("Query: %s\n", query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("Error scanning categories: %v\n", err)
		return
	}
	defer rows.Close()

	categoriesSub := []models.ClothingCategorySub{}
	for rows.Next() {
		var categorySub models.ClothingCategorySub
		if err := rows.Scan(&categorySub.ID, &categorySub.IDClothingCategory, &categorySub.ClothesCatNameSub,
			&categorySub.ClothesCatLocationSub, &categorySub.ClothesPicture1, &categorySub.ClothesPicture2,
			&categorySub.ClothesPicture3, &categorySub.ClothesPicture4, &categorySub.ClothesPicture5,
			&categorySub.ClothesCatStatusSub, &categorySub.CreatedAt, &categorySub.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		categoriesSub = append(categoriesSub, categorySub)
	}
	fmt.Printf("Categories Sub: %v\n", categoriesSub)
	c.JSON(http.StatusOK, categoriesSub)
}

// GetCategorySubByID retrieves a single subcategory by ID
func GetCategorySubByID(c *gin.Context) {
	id := c.Param("id")
	var categorySub models.ClothingCategorySub

	err := db.DB.QueryRow(
		`SELECT id, id_clothing_category, clothes_cat_name_sub, clothes_cat_location_sub, 
         clothes_picture_1, clothes_picture_2, clothes_picture_3, clothes_picture_4, clothes_picture_5, 
         clothes_cat_status_sub, created_at, updated_at FROM clothing_category_sub WHERE id = ?`,
		id,
	).Scan(&categorySub.ID, &categorySub.IDClothingCategory, &categorySub.ClothesCatNameSub,
		&categorySub.ClothesCatLocationSub, &categorySub.ClothesPicture1, &categorySub.ClothesPicture2,
		&categorySub.ClothesPicture3, &categorySub.ClothesPicture4, &categorySub.ClothesPicture5,
		&categorySub.ClothesCatStatusSub, &categorySub.CreatedAt, &categorySub.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subcategory not found"})
		fmt.Printf("Error scanning category sub: %v\n", err)
		return
	}

	c.JSON(http.StatusOK, categorySub)
}

// UpdateCategorySub updates an existing subcategory
func UpdateCategorySub(c *gin.Context) {
	id := c.Param("id")
	var categorySub models.ClothingCategorySub

	if err := c.ShouldBindJSON(&categorySub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	categorySub.UpdatedAt = time.Now()

	_, err := db.DB.Exec(
		`UPDATE clothing_category_sub SET id_clothing_category = ?, clothes_cat_name_sub = ?, 
         clothes_cat_location_sub = ?, clothes_picture_1 = ?, clothes_picture_2 = ?, clothes_picture_3 = ?, 
         clothes_picture_4 = ?, clothes_picture_5 = ?, updated_at = ? WHERE id = ?`,
		categorySub.IDClothingCategory, categorySub.ClothesCatNameSub, categorySub.ClothesCatLocationSub,
		categorySub.ClothesPicture1, categorySub.ClothesPicture2, categorySub.ClothesPicture3,
		categorySub.ClothesPicture4, categorySub.ClothesPicture5, categorySub.UpdatedAt, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("Error updating category sub: %v\n", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subcategory updated successfully"})
}

// DeleteCategorySub soft deletes a subcategory
func DeleteCategorySub(c *gin.Context) {
	id := c.Param("id")

	_, err := db.DB.Exec(
		"UPDATE clothing_category_sub SET clothes_cat_status_sub = 2, updated_at = ? WHERE id = ?",
		time.Now(), id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("Error deleting category sub: %v\n", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subcategory deleted successfully"})
}
