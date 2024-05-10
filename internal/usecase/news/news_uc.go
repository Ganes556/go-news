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
	Update(param ParamUpdate) (err error)
	Delete(param ParamDelete) (err error)
	GetNews(param ParamGetNews) (news []entity.News, err error)
	GetNewsById(ctx context.Context, id uint) (news entity.News, err error)
	GetTotalPost(ctx context.Context) (total int64)
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
		UsersID:      param.UserID,
		CategoriesID: param.Req.CategoryID,
		Title:        param.Req.Title,
		Cover:        url[0],
		Content:      param.Req.Contents,
	}).Error

	if err != nil {
		helper.LogsError(err)
		return err
	}

	return nil
}

type ParamUpdate struct {
	Ctx    context.Context
	Req    req_dto_news.UpdateNews
	UserID uint
}

func (u *ucNews) Update(param ParamUpdate) (err error) {
	err = u.db.WithContext(param.Ctx).Transaction(func(tx *gorm.DB) error {
		var oldData entity.News

		if err2 := tx.First(&oldData, "id = ?", param.Req.ID).Error; err2 != nil {
			if err2 == gorm.ErrRecordNotFound {
				return new(dto_response.Response).Err404("news")
			}
			helper.LogsError(err2)
			return err2
		}
		var newCover string
		if param.Req.Cover != nil {
			err2 := u.Gc.DeleteInStorage(param.Ctx, []string{oldData.Cover})
			if err2 != nil {
				helper.LogsError(err2)
				return err2
			}
			covers, err2 := u.Gc.Upload2Storage(param.Ctx, "cover", []*multipart.FileHeader{param.Req.Cover})
			if err2 != nil {
				helper.LogsError(err2)
				return err2
			}
			newCover = covers[0]
		}
		err2 := tx.Updates(&entity.News{
			Base: entity.Base{
				ID: param.Req.ID,
			},
			UsersID:      param.UserID,
			CategoriesID: param.Req.CategoryID,
			Cover:        newCover,
			Title:        param.Req.Title,
			Content:      param.Req.Contents,
		}).Error

		if err2 != nil {
			helper.LogsError(err2)
			return err2
		}

		return nil
	})
	return	
}

type ParamGetNews struct {
	Ctx   context.Context
	Next  uint
	Limit uint
}

func (u *ucNews) GetNews(param ParamGetNews) (news []entity.News, err error) {
	news = []entity.News{}
	if param.Limit <= 0 {
		param.Limit = 5
	}
	tx := u.db.WithContext(param.Ctx).Omit("content").Order("id ASC").Preload("Users").Preload("Categories")

	if param.Next != 0 {
		tx.Where("id > ?", param.Next)
	}

	err = tx.Limit(int(param.Limit)).
		Find(&news).Error
	return
}
func (u *ucNews) GetNewsById(ctx context.Context, id uint) (news entity.News, err error) {
	news = entity.News{}
	err = u.db.WithContext(ctx).First(&news, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = new(dto_response.Response).Err404("news")
		}
		helper.LogsError(err)
	}
	return
}

func (u *ucNews) GetTotalPost(ctx context.Context) (total int64) {
	u.db.WithContext(ctx).Select("id").Find(&entity.News{}).Count(&total)
	return
}
