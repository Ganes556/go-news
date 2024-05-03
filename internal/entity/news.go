package entity

type UrlImages []string

type News struct {
	Base
	UserID        uint      `json:"-"`
	CategoryID    uint      `json:"-"`
	Category      *Category `json:"-"`
	User          *User     `json:"author"`
	Title         string    `gorm:"type:varchar(150)" json:"title"`
	Cover         string    `gorm:"type:varchar(255)" json:"cover"`
	ContentImages UrlImages `gorm:"serializer:json" json:"-"`
	Content       string    `json:"content"`
	CountView     uint      `json:"count_view"`
}
