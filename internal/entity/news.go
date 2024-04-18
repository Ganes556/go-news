package entity

type News struct {
	Base
	UserID    uint   `json:"-"`
	User      *User  `json:"author"`
	Title     string `gorm:"type:varchar(150)" json:"title"`
	Cover     string `gorm:"type:varchar(255)" json:"Cover"`
	Category  string `gorm:"type:varchar(150);uniqueIndex" json:"category"`
	Content   string `json:"content"`
	CountView uint   `json:"count_view"`
	Date      int64  `gorm:"type:int(11)" json:"date"`
}
