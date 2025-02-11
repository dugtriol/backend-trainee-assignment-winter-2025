package entity

type Transaction struct {
	Id       string `db:"id"`
	FromUser string `db:"from_user"`
	ToUser   string `db:"to_user"`
}
