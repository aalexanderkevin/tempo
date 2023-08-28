package repository

import (
	"context"
	"tempo/model"
)

type News interface {
	Add(ctx context.Context, news *model.News) (*model.News, error)
	Get(ctx context.Context, id *string) (*model.News, error)
	Update(ctx context.Context, id *string, user *model.News) (*model.News, error)
}
