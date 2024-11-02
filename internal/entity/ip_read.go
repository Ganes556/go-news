package entity

import "gorm.io/gorm"

type IpRead struct {
	Base
	IP         string `gorm:"type:varchar(20);uniqueIndex;not null"`
	IpReadable []IpReadable `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

func (i *IpRead) AfterCreate(tx *gorm.DB) (err error) {
	for _, v := range i.IpReadable{
		if err := tx.Table(v.OwnerType).Where("id = ?", v.OwnerID).Update("count_view", gorm.Expr("count_view + 1")).Error; err != nil{
			return err;
		}
		
	}
	return
}