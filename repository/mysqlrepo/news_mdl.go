package mysqlrepo

import (
	"tempo/model"
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type News struct {
	Id          *string
	UserId      *string
	Title       *string
	Description *string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func (n News) FromModel(data model.News) *News {
	return &News{
		Id:          data.Id,
		UserId:      data.UserId,
		Title:       data.Title,
		Description: data.Description,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func (n News) ToModel() *model.News {
	return &model.News{
		Id:          n.Id,
		UserId:      n.UserId,
		Title:       n.Title,
		Description: n.Description,
		CreatedAt:   n.CreatedAt,
		UpdatedAt:   n.UpdatedAt,
	}
}

func (n News) TableName() string {
	return "news"
}

func (n *News) BeforeCreate(db *gorm.DB) error {
	if n.Id == nil {
		db.Statement.SetColumn("id", ksuid.New().String())
	}

	return nil
}
