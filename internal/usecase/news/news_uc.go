package uc_news

import (
	"context"
	"mime/multipart"

	req_dto_news "github.com/news/internal/dto/request/news"
	"github.com/news/internal/entity"
	"github.com/news/pkg"
	"gorm.io/gorm"
)

type NewsUc interface {
}

type newsUc struct {
	Gc pkg.Gcloud
	DB *gorm.DB
}

func NewNewsUc(Gc pkg.Gcloud, DB *gorm.DB) NewsUc {
	return &newsUc{Gc, DB}
}

type ParamCreate struct {
	Ctx    context.Context
	Req    req_dto_news.CreateNews
	UserID uint
}

func (u *newsUc) Create(param ParamCreate) error {
	url, err := u.Gc.Upload2Storge(param.Ctx, "cover", []*multipart.FileHeader{param.Req.Cover})

	if err != nil {
		return err
	}

	err = u.DB.WithContext(param.Ctx).Create(&entity.News{
		UserID: param.UserID,
		Title: param.Req.Title,
		Cover: url[0],
		Category: param.Req.Category,
		Content: param.Req.Contents,
	}).Error
	
	if err != nil {
		return err
	}

	return nil
}
