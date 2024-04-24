package uc_news

import (
	"context"

	"github.com/news/helper"
	dto_error "github.com/news/internal/dto/error"
	req_dto_news "github.com/news/internal/dto/request/news"
	"github.com/news/internal/entity"
	"github.com/news/pkg"
	"gorm.io/gorm"
)

type UcNews interface {
	Create(param ParamCreate) (err error)
	GetNews(param ParamGetNews) (news []entity.News, err error)
	GetNewsById(ctx context.Context, id string) (news entity.News, err error)
}

type ucNews struct {
	Gc pkg.Gcloud
	db *gorm.DB
}

func NewNewsUc(Gc pkg.Gcloud, DB *gorm.DB) UcNews {
	return &ucNews{Gc, DB}
}

type ParamCreate struct {
	Ctx    context.Context
	Req    req_dto_news.CreateNews
	UserID uint
}

func (u *ucNews) Create(param ParamCreate) (err error) {
	// url, err := u.Gc.Upload2Storge(param.Ctx, "cover", []*multipart.FileHeader{param.Req.Cover})

	// if err != nil {
	// 	helper.LogsError(err)
	// 	return err
	// }

	err = u.db.WithContext(param.Ctx).Create(&entity.News{
		UserID:   param.UserID,
		Title:    param.Req.Title,
		Cover:    "",
		Category: param.Req.Category,
		Content:  param.Req.Contents,
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
		param.Limit = 15
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
			err = new(dto_error.ErrResponse).ErrNews404()
		}
	}
	return
}
