package models

import (
	"time"
)

type ClothingCategory struct {
	ID               int       `json:"id"`
	ClothesCatName   string    `json:"clothes_cat_name"`
	ClothesNotes     string    `json:"clothes_notes"`
	ClothesCatStatus int       `json:"clothes_cat_status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ClothingCategorySub struct {
	ID                    int       `json:"id"`
	IDClothingCategory    int       `json:"id_clothing_category"`
	ClothesCatNameSub     string    `json:"clothes_cat_name_sub"`
	ClothesCatLocationSub string    `json:"clothes_cat_location_sub"`
	ClothesPicture1       string    `json:"clothes_picture_1"`
	ClothesPicture2       string    `json:"clothes_picture_2"`
	ClothesPicture3       string    `json:"clothes_picture_3"`
	ClothesPicture4       string    `json:"clothes_picture_4"`
	ClothesPicture5       string    `json:"clothes_picture_5"`
	ClothesCatStatusSub   int       `json:"clothes_cat_status_sub"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type ClothingSize struct {
	ID                    int       `json:"id"`
	IDClothingCategorySub int       `json:"id_clothing_category_sub"`
	ClothesSizeName       string    `json:"clothes_size_name"`
	ClothesSizeNotes      string    `json:"clothes_size_notes"`
	ClothesSizeStatus     int       `json:"clothes_size_status"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type ClothingInventoryMovement struct {
	ID                    int       `json:"id"`
	IDClothingCategory    int       `json:"id_clothing_category"`
	IDClothingSize        int       `json:"id_clothing_size"`
	ClothesMovementAction int       `json:"clothes_movement_action"`
	ClothesQtyIn          int       `json:"clothes_qty_in"`
	ClothesQtyOut         int       `json:"clothes_qty_out"`
	ClothesQtyTotal       int       `json:"clothes_qty_total"`
	ClothesCatStatusSub   int       `json:"clothes_cat_status_sub"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
