//go:build integration
// +build integration

package fixtures

import (
	"backend-trainee-assignment-winter-2025/internal/entity"
	"backend-trainee-assignment-winter-2025/pkg/hasher"
)

type UserBuilder struct {
	instance *entity.User
}

func User() *UserBuilder {
	return &UserBuilder{instance: &entity.User{}}
}

func (b *UserBuilder) ID(value string) *UserBuilder {
	b.instance.Id = value
	return b
}

func (b *UserBuilder) Username(value string) *UserBuilder {
	b.instance.Username = value
	return b
}

func (b *UserBuilder) Password(value string) *UserBuilder {
	var hashPassword string
	hashPassword, _ = hasher.HashPassword(value)
	b.instance.Password = hashPassword
	return b
}

func (b *UserBuilder) Amount(value int) *UserBuilder {
	b.instance.Amount = value
	return b
}

func (b *UserBuilder) P() *entity.User {
	return b.instance
}

func (b *UserBuilder) V() entity.User {
	return *b.instance
}

func (b *UserBuilder) Valid() *UserBuilder {
	return User().Username("user").Password("1234").Amount(1000)
}
