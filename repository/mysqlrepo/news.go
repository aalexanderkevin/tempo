package mysqlrepo

import (
	"context"
	"errors"

	"tempo/model"
	"tempo/repository"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type NewsRepo struct {
	Db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) repository.News {
	return &NewsRepo{
		Db: db,
	}
}

func (u *NewsRepo) Add(ctx context.Context, news *model.News) (*model.News, error) {
	gormModel := News{}.FromModel(*news)

	if err := u.Db.WithContext(ctx).Create(&gormModel).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, model.NewDuplicateError()
		}

		return nil, err
	}

	return gormModel.ToModel(), nil
}

func (u *NewsRepo) Get(ctx context.Context, id *string) (*model.News, error) {
	gormModel := News{}
	q := u.Db.WithContext(ctx).Where("id = ?", *id)

	err := q.First(&gormModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.NewNotFoundError()
		}
		return nil, err
	}

	return gormModel.ToModel(), nil
}

