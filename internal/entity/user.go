package entity

type User struct {
	Base
	Name     string `gorm:"type:varchar(150)" json:"name"`
	Username string `gorm:"type:varchar(150);uniqueIndex" json:"username"`
	Password string `gorm:"type:varchar(255)" json:"-"`
	News     []News `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"news"`
}