package uc_category

import (
	"context"

	"github.com/news/helper"
	dto_response "github.com/news/internal/dto/response"
	"github.com/news/internal/entity"
	"gorm.io/gorm"
)

type UcCategory interface {
	Create(ctx context.Context, category *entity.Category) (err error)
	GetAll(ctx context.Context) (categories []entity.Category, err error)
}

type ucCategory struct {
	db *gorm.DB
}

func NewCategoryUc(DB *gorm.DB) UcCategory {
	return &ucCategory{DB}
}

func (u *ucCategory) Create(ctx context.Context, category *entity.Category) (err error) {
	err = u.db.WithContext(ctx).Create(category).Error
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			err = new(dto_response.Response).ErrCategory409()
		}else {
			helper.LogsError(err)
		}
	}
	return
}

func (u *ucCategory) GetAll(ctx context.Context) (categories []entity.Category, err error){
	categories = []entity.Category{}
	err = u.db.WithContext(ctx).Find(&categories).Error
	if err != nil {
		helper.LogsError(err)
	}
	return
}
