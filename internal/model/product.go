package model

import "time"

type Product struct {
	ID          int64     `db:"id"`
	Category    string    `db:"category"`
	Name        string    `db:"name"`
	Price       float64   `db:"price"`
	Cost        float64   `db:"cost"`
	Description string    `db:"description"`
	Barcode     string    `db:"barcode"`
	ExpiryDate  time.Time `db:"expiry_date"`
	Size        string    `db:"size"`
}
