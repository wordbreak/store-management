package model

type Store struct {
	ID     int64 `db:"id"`
	UserID int64 `db:"user_id"`
}
