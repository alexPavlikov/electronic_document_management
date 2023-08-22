package user

import "context"

type Repository interface {
	SelectUser(ctx context.Context, id string) (User, error)
	SelectUsers(ctx context.Context) ([]User, error)
	SelectUsersByParameter(ctx context.Context) ([]User, error)
	InsertUser(ctx context.Context, user User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id string) error
}
