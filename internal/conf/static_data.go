package conf_internal

import (
	"github.com/news/helper"
	"github.com/news/internal/entity"
)

var DefaultUser = []entity.User{
	// admin
	{
		Base: entity.Base{
			ID: 1,
		},
		Name: "admin",
		Username: "admin",
		Password: func() string{
			e := helper.NewEncryptor()
			phash, _ := e.HashPassword("admin123456")
			return phash
		}(),
	},
}
