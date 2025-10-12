package models

import (
	"time"
)

type ClothingCustomer struct {
	ID          int       `json:"id"`
	CustName    string    `json:"cust_name"`
	CustAddress string    `json:"cust_address"`
	CustCity    string    `json:"cust_city"`
	CustPhone   string    `json:"cust_phone"`
	CustEmail   string    `json:"cust_email"`
	CustNotes   string    `json:"cust_notes"`
	CustStatus  int       `json:"cust_status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
