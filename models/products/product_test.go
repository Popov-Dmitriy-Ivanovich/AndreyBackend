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
		&Category{},
	)

	clearTable[Product](db)
	clearTable[discounts.Discount](db)
	clearTable[productmedia.ProductMedia](db)
	clearTable[Collection](db)
	clearTable[Advert](db)
	clearTable[Category](db)
	discounts := []discounts.Discount{discounts.Discount{NewPrice: 0.69, Style: "Fancy"}}
	productMedia := []productmedia.ProductMedia{productmedia.ProductMedia{File: types.DbFile{""}}}

	product := Product{Name: "test_product", Price: 1.69, IsActive: true, Discounts: discounts, Media: productMedia}
	collection := Collection{Name: "test collection", Products: []*Product{&product}}
	advert := Advert{AdvertText: "test", Products: []*Product{&product}}
	category := Category{Name: "test", Products: []*Product{&product}}
	chldCategory := Category {Name: "test child", ParentCategory: &category, Products: []*Product{&product}}
	productCreateRes := db.Create(&product)
	if productCreateRes.Error != nil {
		return nil, productCreateRes.Error
	}

	collectionCreateRes := db.Create(&collection)
	if collectionCreateRes.Error != nil {
		return nil, collectionCreateRes.Error
	}

	advertCreateRes := db.Create(&advert)
	if advertCreateRes.Error != nil {
		return nil, advertCreateRes.Error
	}

	categoryCreateRes := db.Create(&category)
	if categoryCreateRes.Error != nil {
		return nil, categoryCreateRes.Error
	}

	chldCategoryCreateRes := db.Create(&chldCategory)
	if chldCategoryCreateRes.Error != nil {
		return nil, chldCategoryCreateRes.Error
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

func TestAdvertCascadeDelete(t *testing.T) {
	db, err := InitDb()
	if db == nil || err != nil {
		t.Fatal("CouldNotCreateDB")
	}

	advert := Advert{}
	db.First(&advert)
	foundAdverts := []Advert{}
	db.Table("adverts").Joins("inner join advert_products on advert_products.advert_id = adverts.id").Where("advert_id = ? " , advert.ID).Find(&foundAdverts)
	
	if len(foundAdverts) == 0 {
		t.Error("No products in advert")
	}

	db.Delete(&advert)

	db.Table("adverts").Joins("inner join advert_products on advert_products.advert_id = adverts.id").Where("advert_id = ? " , advert.ID).Find(&foundAdverts)
	
	if len(foundAdverts) != 0 {
		t.Error("advert to product transaction invalid after advert delete")
	}
}

func TestProductWithAdvert(t *testing.T) {
	db, err := InitDb()
	if db == nil || err != nil {
		t.Fatal("CouldNotCreateDB")
	}

	product := Product{}
	db.First(&product)
	foundAdverts := []Advert{}
	
	db.Table("adverts").Joins("inner join advert_products on advert_products.advert_id = adverts.id").Where("product_id = ?", product.ID).Find(&foundAdverts)
	if len(foundAdverts) == 0 {
		t.Error("No adverts found for product")
	}

	db.Delete(&product)

	db.Table("adverts").Joins("inner join advert_products on advert_products.advert_id = adverts.id").Where("product_id = ?", product.ID).Find(&foundAdverts)
	if len(foundAdverts) != 0 {
		t.Error("Advert to product transaction has not been delted after product delete")
	}

}

func TestCategoryDeletition(t *testing.T) {
	db, err := InitDb()
	if db == nil || err != nil {
		t.Fatal("CouldNotCreateDB")
	}

	category := Category{}
	db.Not("parent_id is null").First(&category)
	parrentCategory := Category{}
	db.Find(&parrentCategory, category.ParentID)
	db.Delete(&parrentCategory)

	res := db.Where("category_id = ?", category.ParentID).Find(&[]ProductCategory{})
	if res.RowsAffected != 0 {
		t.Error("product to category transaction does not delete")
	}

	db.First(&category, category.ID)
	if category.ParentID != nil {
		t.Error("child categories are not affected after delete of parrent")
	}
}

func TestProductCategoryDelete(t *testing.T) {
	db, err := InitDb()
	if db == nil || err != nil {
		t.Fatal("CouldNotCreateDB")
	}
	productCategory := ProductCategory{}
	db.First(&productCategory);
	product:= Product{}
	db.First(&product, productCategory.ProductID)
	db.Delete(&product)
	foundProductCategories := []ProductCategory{}
	db.Where("product_id = ?", product.ID).Find(&foundProductCategories)
	if len(foundProductCategories) != 0{
		t.Error("product categories transaction unconsistent after product delete")
	}
}