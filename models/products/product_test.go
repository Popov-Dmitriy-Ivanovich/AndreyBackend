package products

import (
	"fmt"
	"models/discounts"
	"models/productmedia"
	"models/types"
	"testing"

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
		&productmedia.ProductMedia{},
		&Collection{},
	)

	clearTable[Product](db)
	clearTable[discounts.Discount](db)
	clearTable[productmedia.ProductMedia](db)
	clearTable[Collection](db)
	discounts := []discounts.Discount{discounts.Discount{NewPrice: 0.69, Style: "Fancy"}}
	productMedia := []productmedia.ProductMedia{productmedia.ProductMedia{File: types.DbFile{""}}}

	product := Product{Name: "test_product", Price: 1.69, IsActive: true, Discounts: discounts, Media: productMedia}
	collection := Collection{Name: "test collection", Products: []*Product{&product}}

	productCreateRes := db.Create(&product)
	if productCreateRes.Error != nil {
		return nil, productCreateRes.Error
	}

	collectionCreateRes := db.Create(&collection)
	if collectionCreateRes.Error != nil {
		return nil, collectionCreateRes.Error
	}
	// db.Commit()
	return db, nil
}

func TestProductWithDisount(t *testing.T) {
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
	db.Find(&foundDiscounts, &discounts.Discount{ProductID: &product.ID})
	if len(foundDiscounts) != 0 {
		t.Error("Discounts are not deleted after deleted product")
	}
}

func TestProductWithMedia(t *testing.T) {
	db, err := InitDb()
	if db == nil || err != nil {
		t.Fatal("DB is not inited")
	}
	media := productmedia.ProductMedia{}
	db.First(&media)
	product := Product{}
	db.First(&product, media.ProductID)
	db.Delete(&product)
	foundMedia := []productmedia.ProductMedia{}
	db.Find(&foundMedia, &productmedia.ProductMedia{ProductID: product.ID})
	if len(foundMedia) != 0 {
		t.Error("Discounts are not deleted after deleted product")
	}
}

func TestProductWithCollection(t *testing.T) {
	db, err := InitDb()
	if db == nil || err != nil {
		t.Fatal("DB has not been inited")
	}
	product := Product{}
	db.First(&product)
	collections := []Collection{}
	db.Table("collections").Joins("inner join product_collections on product_collections.collection_id = collections.id").Joins("inner join products on product_collections.product_id = products.id").Where("product_id = ?", product.ID).Find(&collections)
	fmt.Println(collections)
	if len(collections) == 0 {
		t.Error("No collections found for product")
	}

	db.Delete(&product)

	db.Table("collections").Joins("inner join product_collections on product_collections.collection_id = collections.id").Joins("inner join products on product_collections.product_id = products.id").Where("product_id = ?", product.ID).Find(&collections)
	if len(collections) != 0 {
		t.Error("Collections found for deleted product")
	}

}

func TestCollectionCascadeDelete(t *testing.T) {
	db, err := InitDb()
	if db == nil || err != nil {
		t.Fatal("DB has not been initialized")
	}

	collection := Collection{}
	db.First(&collection)
	productsFound := []Product{}
	db.Table("products").Joins(
		"inner join product_collections on product_collections.product_id = products.id").Where(
		"collection_id = ?", collection.ID).Find(&productsFound)
	
	if (len(productsFound) == 0 ) {
		t.Error("No products for collection")
	}

	db.Delete(&collection)

	db.Table("products").Joins(
		"inner join product_collections on product_collections.product_id = products.id").Where(
		"collection_id = ?", collection.ID).Find(&productsFound)
	
	if (len(productsFound) != 0 ) {
		t.Error("Products to collections transaction is not updated after collection delete")
	}
}
