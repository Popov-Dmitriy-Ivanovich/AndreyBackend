package products

import (
	// "errors"
	// "models/collections"
	// "log"
	"time"

	"backend/models/discounts"
	"backend/models/productmedia"
	"backend/models/types"

	"gorm.io/gorm"
	// "gorm.io/gorm/clause"
)

// Хоть где-нибудь бы явно это написали в докуметнации: жирным, красным 98 шрифтом
// https://github.com/go-gorm/gorm/issues/6357#issuecomment-1566946772
// https://stackoverflow.com/questions/76762629/how-to-cascade-a-delete-in-gorm
type Product struct {
	gorm.Model
	Name        string
	Desctiption string
	Price       types.PriceType
	IsActive    bool
	Count       uint64
	Picture     types.DbFile
	Discounts   []discounts.Discount `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Media       []productmedia.ProductMedia
	Collection  []*Collection `gorm:"many2many:product_collections"`
	Advert      []*Advert     `gorm:"many2many:advert_products"`
	Categories  []*Category   `gorm:"many2many:product_categories"`
	Reviews 	[]Review
	Carts		[]*Cart		  `gorm:"many2many:product_carts"`
}

type Collection struct {
	gorm.Model
	Name        string
	Description string
	Picture     types.DbFile
	Products    []*Product `gorm:"many2many:product_collections"`
}

type ProductCollection struct {
	ProductID    uint `gorm:"primaryKey"`
	CollectionID uint `gorm:"primaryKey"`
}

type Advert struct {
	gorm.Model
	ExpirationDate time.Time
	AdvertText     string
	Style          string
	Products       []*Product `gorm:"many2many:advert_products"`
}

type AdvertProduct struct {
	ProductID uint `gorm:"primaryKey"`
	AdvertID  uint `gorm:"primaryKey"`
}

type Category struct {
	gorm.Model
	Name           string
	Description    string
	ParentCategory *Category `gorm:"foreignKey:ParentID"`
	ParentID       *uint
	Picture        types.DbFile
	Products       []*Product `gorm:"many2many:product_categories"`
}

type ProductCategory struct {
	ProductID  uint `gorm:"primaryKey"`
	CategoryID uint `gorm:"primaryKey"`
}

type Review struct {
	gorm.Model
	ProductID uint
	UserID uint
	Text string
	Stars uint
}

type ProductCart struct {
	ProductID uint `gorm:"primaryKey"`
	CartID uint `gorm:"primaryKey"`
}

type Cart struct {
	gorm.Model
	UserID uint
	Product []*Product `gorm:"many2many:product_carts"`
}

func (category *Category) AfterDelete(db *gorm.DB) error {
	res := db.Where("category_id = ?", category.ID).Delete(&ProductCategory{})
	if res.Error != nil {
		return res.Error
	}
	res = db.Model(&Category{}).Where("parent_id = ?", category.ID).Update("parent_id", nil)
	return res.Error
}

func (prod *Product) AfterDelete(db *gorm.DB) error {
	//return nil

	res := db.Where("product_id = ?", prod.ID).Delete(&discounts.Discount{})
	if res.Error != nil {
		return res.Error
	}

	res = db.Where("product_id = ?", prod.ID).Delete(&productmedia.ProductMedia{})
	if res.Error != nil {
		return res.Error
	}
	res = db.Where("product_id = ?", prod.ID).Delete(&ProductCollection{})
	if res.Error != nil {
		return res.Error
	}

	res = db.Where("product_id = ?", prod.ID).Delete(&AdvertProduct{})
	if res.Error != nil {
		return res.Error
	}
	res = db.Where("product_id = ?", prod.ID).Delete(&ProductCategory{})

	if res.Error != nil {
		return res.Error
	}

	res = db.Where("product_id = ?", prod.ID).Delete(&Review{})
	if res.Error != nil {
		return res.Error
	}
	
	res = db.Where("product_id = ?", prod.ID).Delete(&ProductCart{})
	// res = db.Delete(&ProductCollection{}, &ProductCollection{ProductID: prod.ID})

	return res.Error
}

func (collection *Collection) AfterDelete(db *gorm.DB) error {
	res := db.Where("collection_id = ?", collection.ID).Delete(&ProductCollection{})
	return res.Error
}

func (advert *Advert) AfterDelete(db *gorm.DB) error {
	res := db.Where("advert_id = ?", advert.ID).Delete(&AdvertProduct{})
	return res.Error
}

func (cart *Cart) AfterDelete (db *gorm.DB) error {
	res := db.Where("cart_id = ?", cart.ID).Delete(&ProductCart{})
	return res.Error
}