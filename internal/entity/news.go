package entity

type News struct {
	Base
	UserID    uint
	Title     string `gorm:"type:varchar(150)" json:"title"`
	Image     string `gomr:"type:varchar(255)" json:"image"`
	Category  string `gorm:"type:varchar(150);uniqueIndex" json:"category"`
	Content   string `json:"content"`
	CountView uint   `json:"count_view"`
	Date      int64  `gorm:"type:int(11)" json:"date"`
}
