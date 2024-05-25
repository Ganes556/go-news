package entity

type IpRead struct {
	Base
	IP   string `gorm:"type:varchar(20);uniqueIndex;not null"`
	News []News `gorm:"constraint:OnDelete:CASCADE;many2many:ipread_news"`
}
