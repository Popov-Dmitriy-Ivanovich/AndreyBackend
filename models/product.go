package models
import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string
	Desctiption string
	Price       PriceType
	IsActive    bool
	Count       uint64
	Picture     DbFile
	Discounts 	[]Discount  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
