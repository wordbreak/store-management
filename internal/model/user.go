package model

type User struct {
	Id          int64  `db:"id"`
	PhoneNumber string `db:"phone_number"`
	Password    string `db:"password"`
}
