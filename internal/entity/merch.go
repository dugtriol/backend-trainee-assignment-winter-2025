package entity

type Merch struct {
	Id    string `db:"id"`
	Name  string `db:"name"`
	Price int    `db:"price"`
}
