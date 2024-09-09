package products

import (
	"testing"
	"models/discounts"
	// "models/products"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func clearTable[T any](db *gorm.DB) error {
	var rows []T
	res := db.Find(&rows)
	if res.Error != nil {
		return res.Error
	}
	for _, row := range rows {
		delRes := db.Unscoped().Delete(&row)
		if delRes.Error != nil {
			return delRes.Error
		}
	}

	return nil
}

func InitDb() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		//panic("failed to connect database")
		return nil, err
	}
	db.AutoMigrate(
		&Product{},
		&discounts.Discount{},
	)

	clearTable[Product](db)
	clearTable[discounts.Discount](db)
	
	discounts := []discounts.Discount{discounts.Discount{NewPrice: 0.69, Style: "Fancy"}}
	product := Product{Name: "test_product", Price: 1.69, IsActive: true, Discounts: discounts}

	productCreateRes := db.Create(&product)
	if (productCreateRes.Error != nil){
		return nil, productCreateRes.Error
	}

	// db.Commit()
	return db, nil
}

func TestProductWithDisount (t *testing.T) {
	db, err := InitDb()
	if db == nil || err != nil {
		t.Fatal("DB is not inited")
	}
	discount := discounts.Discount{}
	db.First(&discount)
	product := Product{}
	db.First(&product, discount.ProductID)
	db.Delete(&product)
	foundDiscounts := []discounts.Discount{}
	db.Find(&foundDiscounts,&discounts.Discount{ProductID: &product.ID})
	if (len(foundDiscounts) != 0) {
		t.Error("Discounts are not deleted after deleted product")
	}
}