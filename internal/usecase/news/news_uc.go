package uc_news

import (
	"context"
	"mime/multipart"

	"github.com/news/helper"
	req_dto_news "github.com/news/internal/dto/request/news"
	dto_response "github.com/news/internal/dto/response"
	"github.com/news/internal/entity"
	"github.com/news/pkg"
	"gorm.io/gorm"
)

type UcNews interface {
	Create(param ParamCreate) (err error)
	Delete(param ParamDelete) (err error)
	GetNews(param ParamGetNews) (news []entity.News, err error)
	GetDistinctCategory(ctx context.Context) (categories []string, err error)
	GetNewsById(ctx context.Context, id string) (news entity.News, err error)
}

type ucNews struct {
	Gc pkg.Gcloud
	db *gorm.DB
}

func NewNewsUc(Gc pkg.Gcloud, DB *gorm.DB) UcNews {
	return &ucNews{Gc, DB}
}

type ParamDelete struct {
	Ctx context.Context
	Req req_dto_news.DeleteNews
}

func (u *ucNews) Delete(param ParamDelete) (err error) {
	err = u.db.WithContext(param.Ctx).Transaction(func(tx *gorm.DB) error {
		oldNews := new(entity.News)
		if err := tx.First(oldNews, "id = ?", param.Req.ID).Error; err != nil {
			return err
		}
		if err := u.Gc.DeleteInStorage(param.Ctx, []string{oldNews.Cover}); err != nil {
			return err
		}
		if err := tx.Delete(oldNews, "id = ?", param.Req.ID).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		helper.LogsError(err)
	}
	return
}

type ParamCreate struct {
	Ctx    context.Context
	Req    req_dto_news.CreateNews
	UserID uint
}

func (u *ucNews) Create(param ParamCreate) (err error) {
	url, err := u.Gc.Upload2Storage(param.Ctx, "cover", []*multipart.FileHeader{param.Req.Cover})

	if err != nil {
		helper.LogsError(err)
		return err
	}

	err = u.db.WithContext(param.Ctx).Create(&entity.News{
		UserID:     param.UserID,
		CategoryID: param.Req.CategoriID,
		Title:      param.Req.Title,
		Cover:      url[0],
		Content:    param.Req.Contents,
	}).Error

	if err != nil {
		helper.LogsError(err)
		return err
	}

	return nil
}

type ParamGetNews struct {
	Ctx   context.Context
	Next  uint
	Limit uint
}

func (u *ucNews) GetNews(param ParamGetNews) (news []entity.News, err error) {
	news = []entity.News{}
	if param.Limit <= 0 {
		param.Limit = 10
	}
	tx := u.db.WithContext(param.Ctx).Order("id ASC").Preload("User")

	if param.Next != 0 {
		tx.Where("id > ?", param.Next)
	}

	err = tx.Limit(int(param.Limit)).
		Find(&news).Error

	return
}
func (u *ucNews) GetNewsById(ctx context.Context, id string) (news entity.News, err error) {
	news = entity.News{}
	err = u.db.WithContext(ctx).First(&news, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = new(dto_response.Response).ErrNews404()
		}
		helper.LogsError(err)
	}
	return
}

func (u *ucNews) GetDistinctCategory(ctx context.Context) (categories []string, err error) {
	news := new(entity.News)
	err = u.db.WithContext(ctx).Model(news).Select("category").Distinct("category").Pluck("category", &categories).Error
	if err != nil {
		helper.LogsError(err)
	}
	return
}
