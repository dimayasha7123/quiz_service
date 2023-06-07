package domain

import "context"

type RepoUser struct {
	ID   int64
	QSID int64
	Name string
}

type repository interface {
	GetAllUsers(ctx context.Context) ([]RepoUser, error)
	AddUser(ctx context.Context, user RepoUser) error
	GetUserByID(ctx context.Context, id int64) (RepoUser, error)
}
