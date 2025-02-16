package entity

type Inventory struct {
	ID         string `db:"id"`
	CustomerID string `db:"customer_id"`
	Type       string `db:"type"`
	Quantity   int    `db:"quantity"`
}
