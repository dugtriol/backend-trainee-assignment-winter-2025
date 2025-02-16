//go:build integration
// +build integration

package fixtures

import (
	"backend-trainee-assignment-winter-2025/internal/entity"
)

type TransactionBuilder struct {
	instance *entity.Transaction
}

func Transaction() *TransactionBuilder {
	return &TransactionBuilder{instance: &entity.Transaction{}}
}

func (b *TransactionBuilder) FromUser(value string) *TransactionBuilder {
	b.instance.FromUser = value
	return b
}

func (b *TransactionBuilder) ToUser(value string) *TransactionBuilder {
	b.instance.ToUser = value
	return b
}

func (b *TransactionBuilder) Amount(value int) *TransactionBuilder {
	b.instance.Amount = value
	return b
}

func (b *TransactionBuilder) P() *entity.Transaction {
	return b.instance
}

func (b *TransactionBuilder) V() entity.Transaction {
	return *b.instance
}
