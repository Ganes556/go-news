package cron_func

import (
	"github.com/news/internal/entity"
	"gorm.io/gorm"
)

type CronFunc interface {
	DBResetIpread()
}

type cronFunc struct {
	DB *gorm.DB
}

func NewCronFunc(DB *gorm.DB) CronFunc{
	return &cronFunc{DB}
}

func (c *cronFunc) DBResetIpread() {
	err := c.DB.Delete(&entity.IpRead{}, "1 = 1").Association("IpReadable").Clear().Error
	if err != nil {
		panic(err)
	}
}