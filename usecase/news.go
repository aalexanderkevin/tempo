package usecase

import (
	"context"

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

