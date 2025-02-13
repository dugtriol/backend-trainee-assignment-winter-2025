package entity

type Inventory struct {
	Id         string `db:"id"`
	CustomerId string `db:"customer_id"`
	Type       string `db:"type"`
	Quantity   int    `db:"quantity"`
}
