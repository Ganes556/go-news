package uc_user

import (
	"context"

	"github.com/news/helper"
	req_dto_user "github.com/news/internal/dto/request/user"
	dto_response "github.com/news/internal/dto/response"
	"github.com/news/internal/entity"
	"gorm.io/gorm"
)

type UcUser interface {
	Login(ctx context.Context, req req_dto_user.Login) (user *entity.Users, err error)
	EditProfile(ctx context.Context, id uint, req req_dto_user.EditProfile) (err error)
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

func (u *ucUser) Login(ctx context.Context, req req_dto_user.Login) (user *entity.Users, err error) {
	user = &entity.Users{}
	err = u.db.WithContext(ctx).First(user,"username = ?", req.Username).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = new(dto_response.Response).ErrAuth400()
			return
		}		
	}

	if !u.encryptor.CheckPasswordHash(req.Password, user.Password) {
		err = new(dto_response.Response).ErrAuth400()
		return
	}

	return
}

func (u *ucUser) EditProfile(ctx context.Context, id uint, req req_dto_user.EditProfile) (err error) {
	err = u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		user := new(entity.Users)
		errTx := tx.Where("id = ?", id).First(user).Error
		if errTx != nil {
			return errTx
		}
		var newHaspassword string
		if req.OldPassword != "" {
			check := u.encryptor.CheckPasswordHash(req.OldPassword, user.Password)
			if !check {
				return new(dto_response.Response).ErrChangePassword400()
			}
			newHash, err := u.encryptor.HashPassword(req.NewPassword)
			if err != nil {
				return new(dto_response.Response).ErrAuth500()
			}
			newHaspassword = newHash
		}

		errTx = tx.Model(user).Updates(&entity.Users{
			Name: req.Name,
			Username: req.Username,
			Password: newHaspassword,
		}).Error

		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = new(dto_response.Response).ErrAuth400()
			return
		}		
	}
	return
}