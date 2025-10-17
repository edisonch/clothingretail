package models

import (
	"time"
)

type ClothingRental struct {
	ID                          int       `json:"id"`
	IDClothingCategorySub       int       `json:"id_clothing_category_sub"`
	IDClothingSize              int       `json:"id_clothing_size"`
	IDClothingCustomer          int       `json:"id_clothing_customer"`
	ClothesQtyRent              int       `json:"clothes_qty_rent"`
	ClothesQtyReturn            int       `json:"clothes_qty_return"`
	ClothesRentDateBegin        time.Time `json:"clothes_rent_date_begin"`
	ClothesRentDateEnd          time.Time `json:"clothes_rent_date_end"`
	ClothesRentDateActualPickup time.Time `json:"clothes_rent_date_actual_pickup"`
	ClothesRentDateActualReturn time.Time `json:"clothes_rent_date_actual_return"`
	ClothesRentStatus           int       `json:"clothes_rent_status"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
}

type RentalRequest struct {
	IDClothingCategorySub int    `json:"id_clothing_category_sub" binding:"required"`
	IDClothingSize        int    `json:"id_clothing_size" binding:"required"`
	IDClothingCustomer    int    `json:"id_clothing_customer" binding:"required"`
	ClothesQtyRent        int    `json:"clothes_qty_rent" binding:"required"`
	RentDateBegin         string `json:"rent_date_begin" binding:"required"`
	RentDateEnd           string `json:"rent_date_end" binding:"required"`
}

type ReturnRequest struct {
	RentalID         int `json:"rental_id" binding:"required"`
	ClothesQtyReturn int `json:"clothes_qty_return" binding:"required"`
}
