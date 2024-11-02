package entity

type UrlImages []string

type News struct {
	Base
	UsersID      uint        `gorm:"not null"`
	CategoriesID uint        `gorm:"not null"`
	Categories   *Categories `json:"-"`
	Users        *Users
	Slug         string `gorm:"type:varchar(150);uniqueIndex;not null"`
	Title        string `gorm:"type:varchar(150);uniqueIndex;not null"`
	Cover        string `gorm:"type:varchar(255);not null"`
	Content      string `gorm:"not null"`
	CountView    uint
	IpReadable   []IpReadable `gorm:"polymorphic:Owner"`
}
