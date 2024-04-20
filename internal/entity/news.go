package entity

type UrlImages []string

type News struct {
	Base
	UserID        uint      `json:"-"`
	User          *User     `json:"author"`
	Title         string    `gorm:"type:varchar(150)" json:"title"`
	Cover         string    `gorm:"type:varchar(255)" json:"Cover"`
	Category      string    `gorm:"type:varchar(150);uniqueIndex" json:"category"`
	ContentImages UrlImages `gorm:"type:serializer:json" json:"-"`
	Content       string    `json:"content"`
	CountView     uint      `json:"count_view"`
	Date          int64     `gorm:"type:int(11)" json:"date"`
}
