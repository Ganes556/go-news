package entity

type Category struct {
	Base
	Name string `gorm:"type:varchar(150);uniqueIndex" json:"name"`
	News []News `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"news"`
}
