package auth

import (
	"context"
	"errors"
)

type contextKey int

const (
	customerContextKey contextKey = iota
)

var ErrNoCustomerInContext = errors.New("no customer in context")

type Customer struct {
	UUID  string
	Name  string
	Role  string
	Email string
}

func CustomerFromContext(ctx context.Context) (Customer, error) {
	customer, ok := ctx.Value(customerContextKey).(Customer)
	if ok {
		return customer, nil
	}

	return Customer{}, ErrNoCustomerInContext
}
