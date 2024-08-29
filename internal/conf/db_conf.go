package conf_internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	sqlNative "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	fmysql "github.com/gofiber/storage/mysql"
	fsqlite3 "github.com/gofiber/storage/sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/news/internal/entity"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
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

func NewGorm(param ParamNewGorm) (*gorm.DB, fiber.Storage) {
	var db *gorm.DB
	var err error
	var dsn string
	var storage fiber.Storage
	if os.Getenv("DB_CONNECTION") == "sqlite" {
		dir := "db/sqlite"
		if err := os.MkdirAll(dir, 0775); err != nil {
			panic(err)
		}
		dsn = dir + "/" + os.Getenv("DB_DATABASE") + ".db"
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})

		storage = fsqlite3.New(fsqlite3.Config{
			Database:        dsn,
			Table:           "fiber_storage",
			Reset:           false,
			GCInterval:      10 * time.Second,
			MaxOpenConns:    100,
			MaxIdleConns:    100,
			ConnMaxLifetime: 1 * time.Second,
		})

	} else {
		if os.Getenv("ENVIRONMENT") == "PRODUCTION" {
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				param.Username, param.Password, param.Host, param.Port, param.Database)
		} else {
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				"root", "", "localhost", "3306", param.Database)
		}
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		storage = fmysql.New(fmysql.Config{
			Host: param.Host,
			Port: func() int {
				port, _ := strconv.Atoi(param.Port) // Handle error appropriately
				return port
			}(),
			Database:   param.Database,
			Username:   param.Username,
			Password:   param.Password,
			Table:      "fiber_store",
			Reset:      false,
			GCInterval: 10 * time.Second,
		})
	}
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(new(entity.Users), new(entity.Categories), new(entity.News), new(entity.IpRead)); err != nil {
		panic(err)
	}
	
	if err := addDefultValues(db, DefaultUser); err != nil {
		panic(err)
	}

	return db, storage
}

func addDefultValues(db *gorm.DB, values ...interface{}) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, v := range values {
			if err := tx.Create(v).Error; err != nil {
				if mysqlErr, ok := err.(*sqlNative.MySQLError); ok && mysqlErr.Number == 1062 {
					// Ignore the duplicate entry error
					continue
				}
				if strings.Contains(err.Error(), "UNIQUE constraint failed") {
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
