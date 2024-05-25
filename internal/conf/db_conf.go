package conf_internal

import (
	"fmt"

	sqlNative "github.com/go-sql-driver/mysql"
	"github.com/news/internal/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type ParamNewGorm struct {
	Username string
	Password string
	Database string
	Port     string
	Host     string
}

func NewGorm(param ParamNewGorm) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", param.Username, param.Password, param.Host, param.Port, param.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(new(entity.Users), new(entity.Categories), new(entity.News), new(entity.IpRead)); err != nil {
		panic(err)
	}

	if err := addDefultValues(db, DefaultUser); err != nil {
		panic(err)
	}

	return db
}

func addDefultValues(db *gorm.DB, values ...interface{}) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, v := range values {
			if err := tx.Create(v).Error; err != nil {
				if mysqlErr, ok := err.(*sqlNative.MySQLError); ok && mysqlErr.Number == 1062 {
					// Ignore the duplicate entry error
					continue
				}
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
