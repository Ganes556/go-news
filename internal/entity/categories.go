package entity

type Categories struct {
	Base
	Name       string       `gorm:"type:varchar(150);uniqueIndex" json:"name"`
	News       []News       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"news"`
	IpReadable []IpReadable `gorm:"polymorphic:Owner"`
	CountView  uint
}
