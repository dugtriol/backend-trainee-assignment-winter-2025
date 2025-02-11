package entity

type Inventory struct {
	Id         string `db:"id"`
	Type       string `db:"type"`
	Quantity   int    `db:"quantity"`
	CustomerId string `db:"customer_id"`
}
