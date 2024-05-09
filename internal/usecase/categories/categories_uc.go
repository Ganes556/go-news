package uc_categories

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/news/helper"
	req_dto_categories "github.com/news/internal/dto/request/categories"
	dto_response "github.com/news/internal/dto/response"
	"github.com/news/internal/entity"
	"github.com/news/pkg"
	"gorm.io/gorm"
)

type UcCategories interface {
	Create(ctx context.Context, req req_dto_categories.CreateCategory) (err error)
	Update(ctx context.Context, req req_dto_categories.UpdateCategory) (err error)
	GetAll(ctx context.Context) (categories []entity.Categories, err error)
	Delete(ctx context.Context, id uint) (err error)
}

type ucCategories struct {
	db *gorm.DB
	gc pkg.Gcloud
}

func NewCategoriesUc(DB *gorm.DB, gc pkg.Gcloud) UcCategories {
	return &ucCategories{DB, gc}
}

func (u *ucCategories) Create(ctx context.Context, req req_dto_categories.CreateCategory) (err error) {
	newCategory := entity.Categories{
		Name: req.Name,
	}
	err = u.db.WithContext(ctx).Create(&newCategory).Error
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			err = new(dto_response.Response).Err409("category")
		} else {
			helper.LogsError(err)
		}
	}
	return
}

func (u *ucCategories) Update(ctx context.Context, req req_dto_categories.UpdateCategory) (err error) {
	newCategory := entity.Categories{
		Base: entity.Base{
			ID: req.ID,
		},
		Name: req.Name,
	}
	fmt.Println("request", req)
	err = u.db.WithContext(ctx).Updates(&newCategory).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if ok := errors.As(err, &mysqlErr); ok {
			if mysqlErr.Number == 1062 { // duplicate entries based on the index rule
				err = new(dto_response.Response).Err409("category")
			}
		} else {
			helper.LogsError(err)
		}
	}
	return
}

func (u *ucCategories) GetAll(ctx context.Context) (categories []entity.Categories, err error) {
	categories = []entity.Categories{}
	err = u.db.WithContext(ctx).Find(&categories).Error
	if err != nil {
		helper.LogsError(err)
	}
	return
}

func (u *ucCategories) Delete(ctx context.Context, id uint) (err error) {
	err = u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		category := &entity.Categories{}
		err2 := tx.Preload("News").First(category,"id = ?", id).Error	
		if err2 != nil {
			if err2 == gorm.ErrRecordNotFound {
				err2 = new(dto_response.Response).Err404("category")
				return err2
			}
			helper.LogsError(err2)
			return err2
		}
		err2 = u.db.WithContext(ctx).Delete(category).Error
		if err != nil {		
			helper.LogsError(err2)
			return err
		}

		if category.News != nil {
			var objNames = make([]string, len(category.News))
			for i, v := range category.News {
				objNames[i] = v.Cover
			}
			err2 = u.gc.DeleteInStorage(ctx, objNames)
			if err2 != nil {
				return err2
			}
		}
		return nil
	})
	
	return
}
