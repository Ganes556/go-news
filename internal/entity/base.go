package entity

type Base struct {
	ID        uint  `gorm:"primarykey" json:"id"`
	CreatedAt int64 `gorm:"autoCreateTime" json:"created_at"`
}
