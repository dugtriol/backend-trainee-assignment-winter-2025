package entity

type User struct {
	ID       string `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Amount   int    `db:"amount"`
}
