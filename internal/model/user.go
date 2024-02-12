package model

type User struct {
	ID          int64  `db:"id"`
	PhoneNumber string `db:"phone_number"`
	Password    string `db:"password"`
}
