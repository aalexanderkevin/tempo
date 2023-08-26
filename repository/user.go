package repository

import (
	"context"
	"tempo/model"
)

type User interface {
	Add(ctx context.Context, user *model.User) (*model.User, error)
	Get(ctx context.Context, filter UserGetFilter) (*model.User, error)
	Update(ctx context.Context, id string, user *model.User) (*model.User, error)
}

type UserGetFilter struct {
	Id    *string
	Email *string
}
