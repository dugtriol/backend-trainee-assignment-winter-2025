package entity

type Transaction struct {
	ID       string `db:"id"`
	FromUser string `db:"from_user"`
	ToUser   string `db:"to_user"`
	Amount   int    `db:"amount"`
}
