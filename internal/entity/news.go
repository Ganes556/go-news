package entity

type UrlImages []string

type News struct {
	Base
	UserID        uint      `json:"-"`
	User          *User     `json:"author"`
	Title         string    `gorm:"type:varchar(150)" json:"title"`
	Cover         string    `gorm:"type:varchar(255)" json:"Cover"`
	Category      string    `gorm:"type:varchar(150);index" json:"category"`
	ContentImages UrlImages `gorm:"serializer:json" json:"-"`
	Content       string    `json:"content"`
	CountView     uint      `json:"count_view"`
}
