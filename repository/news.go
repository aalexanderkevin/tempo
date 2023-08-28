package repository

import (
	"context"
	"tempo/model"
)

type News interface {
	Add(ctx context.Context, news *model.News) (*model.News, error)
}
