package handlers

import (
	"clothingretail/db"
	"clothingretail/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateCategory handles creating a new clothing category
func CreateCategory(c *gin.Context) {
	var category models.ClothingCategory

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	category.ClothesCatStatus = 1

	result, err := db.DB.Exec(
		"INSERT INTO clothing_category (clothes_cat_name, clothes_notes, clothes_cat_status, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		category.ClothesCatName, category.ClothesNotes, category.ClothesCatStatus, category.CreatedAt, category.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	category.ID = int(id)

	c.JSON(http.StatusCreated, category)
}

// GetCategories retrieves all clothing categories
func GetCategories(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, clothes_cat_name, clothes_notes, clothes_cat_status, created_at, updated_at FROM clothing_category WHERE clothes_cat_status = 1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var categories []models.ClothingCategory
	for rows.Next() {
		var category models.ClothingCategory
		if err := rows.Scan(&category.ID, &category.ClothesCatName, &category.ClothesNotes, &category.ClothesCatStatus, &category.CreatedAt, &category.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoryByID retrieves a single category by ID
func GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	var category models.ClothingCategory

	err := db.DB.QueryRow(
		"SELECT id, clothes_cat_name, clothes_notes, clothes_cat_status, created_at, updated_at FROM clothing_category WHERE id = ?",
		id,
	).Scan(&category.ID, &category.ClothesCatName, &category.ClothesNotes, &category.ClothesCatStatus, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateCategory updates an existing category
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.ClothingCategory

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.UpdatedAt = time.Now()

	_, err := db.DB.Exec(
		"UPDATE clothing_category SET clothes_cat_name = ?, clothes_notes = ?, updated_at = ? WHERE id = ?",
		category.ClothesCatName, category.ClothesNotes, category.UpdatedAt, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

// DeleteCategory soft deletes a category
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	_, err := db.DB.Exec(
		"UPDATE clothing_category SET clothes_cat_status = 2, updated_at = ? WHERE id = ?",
		time.Now(), id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
