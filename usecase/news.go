package usecase

import (
	"context"
	"errors"

	"tempo/container"
	"tempo/helper"
	"tempo/model"
	"tempo/repository"
)

type News struct {
	repository.News
}

func NewNews(n *container.Container) *News {
	return &News{
		News: n.NewsRepo(),
	}
}

func (n *News) Add(ctx context.Context, req *model.News) (*model.News, error) {
	logger := helper.GetLogger(ctx).WithField("method", "usecase.News.Add")

	if err := req.Validate(); err != nil {
		logger.WithError(err).Warning("Not Valid Request")
		return nil, model.NewParameterError(helper.Pointer(err.Error()))
	}

	res, err := n.News.Add(ctx, req)
	if err != nil {
		logger.WithError(err).Warning("Failed insert News")
		return nil, err
	}

	return res, nil
}

func (n *News) Get(ctx context.Context, id *string) (*model.News, error) {
	logger := helper.GetLogger(ctx).WithField("method", "usecase.News.Get")

	if id == nil {
		err := errors.New("id is missing")
		logger.WithError(err).Warning("Not Valid Request")
		return nil, model.NewParameterError(helper.Pointer(err.Error()))
	}

	user, err := n.News.Get(ctx, id)
	if err != nil {
		logger.WithError(err).Warning("Failed get News")
		return nil, err
	}

	return user, nil
}

