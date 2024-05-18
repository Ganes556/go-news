package entity

type UrlImages []string

type News struct {
	Base
	UsersID       uint        `gorm:"not null"`
	CategoriesID  uint        `gorm:"not null"`
	Categories    *Categories `json:"-"`
	Users         *Users
	Title         string    `gorm:"type:varchar(150);uniqueIndex;not null"`
	Cover         string    `gorm:"type:varchar(255);not null"`
	ContentImages UrlImages `gorm:"serializer:json"`
	Content       string    `gorm:"not null"`
	CountView     uint
}
