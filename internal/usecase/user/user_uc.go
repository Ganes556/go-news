package uc_user

import (
	"context"

	"github.com/news/helper"
	dto_error "github.com/news/internal/dto/error"
	req_dto_user "github.com/news/internal/dto/request/user"
	"github.com/news/internal/entity"
	"gorm.io/gorm"
)

type UcUser interface {
	Login(ctx context.Context, req req_dto_user.Login) (user *entity.User, err error)
	GetNews(ctx context.Context) (news []entity.News, err error)
	GetNewsById(ctx context.Context, id string) (news entity.News, err error)
}

type ucUser struct {
	db        *gorm.DB
	encryptor helper.Encryptor
}

func NewUcUser(db *gorm.DB, encryptor helper.Encryptor) UcUser {
	return &ucUser{
		db:        db,
		encryptor: encryptor,
	}
}

func (u *ucUser) Login(ctx context.Context, req req_dto_user.Login) (user *entity.User, err error) {
	user = &entity.User{}
	err = u.db.WithContext(ctx).First(user,"username = ?", req.Username).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = new(dto_error.ErrResponse).ErrAuth400()
			return
		}		
	}
	return
}
func (u *ucUser) GetNews(ctx context.Context) (news []entity.News, err error) {
	news = []entity.News{}
	err = u.db.WithContext(ctx).Order("id ASC").Find(&news).Error
	return
}
func (u *ucUser) GetNewsById(ctx context.Context, id string) (news entity.News, err error) {
	news = entity.News{}
	err = u.db.WithContext(ctx).First(&news, "id = ?", id ).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = new(dto_error.ErrResponse).ErrNews404()
		}
	}
	return
}
