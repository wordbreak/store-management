package model

import "time"

type Product struct {
	ID           int64     `db:"id" json:"id"`
	Category     string    `db:"category" json:"category"`
	Name         string    `db:"name" json:"name"`
	AbstractName string    `db:"abstract_name" json:"-"`
	Price        float64   `db:"price" json:"price"`
	Cost         float64   `db:"cost" json:"cost"`
	Description  string    `db:"description" json:"description"`
	Barcode      string    `db:"barcode" json:"barcode"`
	ExpiryDate   time.Time `db:"expiry_date" json:"expiryDate"`
	Size         string    `db:"size" json:"size"`
}
