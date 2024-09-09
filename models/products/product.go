package products

import (
	// "errors"
	"models/types"
	"models/discounts"
	"gorm.io/gorm"
	// "gorm.io/gorm/clause"
)

// Хоть где-нибудь бы явно это написали в докуметнации: жирным, красным 98 шрифтом
//https://github.com/go-gorm/gorm/issues/6357#issuecomment-1566946772
//https://stackoverflow.com/questions/76762629/how-to-cascade-a-delete-in-gorm
type Product struct {
	gorm.Model
	Name        string
	Desctiption string
	Price       types.PriceType
	IsActive    bool
	Count       uint64
	Picture     types.DbFile
	Discounts 	[]discounts.Discount  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (prod *Product) AfterDelete (db *gorm.DB) error {
	// discounts := []Discount{}
	//db.Where(&Discount{ProductID: &prod.ID}, "product_id").Delete(&Discount{})
	// db.Delete(&discounts)
	//if prod.ID == 0 {
	//	return errors.New("wrong id in delete hook")
	//}
	res := db.Where("product_id = ?", prod.ID).Delete(&discounts.Discount{})
	return res.Error
}