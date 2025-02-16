package entity

type Merch struct {
	ID    string `db:"id"`
	Name  string `db:"name"`
	Price int    `db:"price"`
}
