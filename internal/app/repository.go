package app

import "context"

type Repository interface {
	AddUser(context.Context, string) (int64, error)
}
