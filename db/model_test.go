package db

import (
	// "fmt"
	_ "fmt"
	"testing"

	// "gorm.io/gen"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	// "gorm.io/clause"
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
	db.AutoMigrate(&Product{},
		&Discount{},
		&ProductMedia{},
		&Category{},
		&Advert{},
		&User{},
		&ArticleMedia{},
		&Article{},
		&Collection{},
		&Cart{},
		&NewsSubscribers{},
		&ProductReview{},
		&CartProduct{},
		&ProductCategory{},
		&ProductAdvert{},
		&ProductCollection{})
	
	clearTable[Product](db)
	clearTable[Discount](db)
	clearTable[ProductMedia](db)
	clearTable[Category](db)
	clearTable[Advert](db)
	clearTable[User](db)
	clearTable[ArticleMedia](db)
	clearTable[Article](db)
	clearTable[Collection](db)
	clearTable[Cart](db)
	clearTable[NewsSubscribers](db)
	clearTable[ProductReview](db)
	clearTable[CartProduct](db)
	clearTable[ProductCategory](db)
	clearTable[ProductAdvert](db)
	clearTable[ProductCollection](db)
	
	product := &Product{Name: "test_product",
		Desctiption: "test_product test description",
		Price:       0.69,
		IsActive:    true,
		Count:       248432,
		Picture:     DbFile{""}}
	db.Create(product)

	discount := &Discount{Product: *product, NewPrice: 0.34, Style: "fancy"}
	db.Create(discount)
	
	productMedia := &ProductMedia{Product: *product, File: DbFile{""}}
	db.Create(productMedia)
	
	category := &Category{Name: "cats", Description: "desc", Picture: DbFile{""}}
	db.Create(category)

	chldCategory := &Category{Name: "cute cats", Description: "desc", ParrentCategory: category}
	db.Create(chldCategory)

	advert := &Advert{Text: "advert text"}
	db.Create(advert)

	user := &User{Email: "huy"}
	db.Create(user)

	article := &Article{Author: *user, Html: DbFile{""}}
	db.Create(article)

	articleMedia := &ArticleMedia{File: DbFile{""}, Article: *article}
	db.Create(articleMedia)	
	
	collection := &Collection{Name: "test collection", Description: "desc", Picture: DbFile{""}}
	db.Create(collection)

	cart := &Cart{User: *user}
	db.Create(cart)

	newsSubscribers := &NewsSubscribers{Email: "huy"}
	db.Create(newsSubscribers)

	productReview := &ProductReview { Product: *product, User: *user, Text: "text", Stars: 8}
	db.Create(productReview)

	cartProduct := &CartProduct{Cart: *cart, Product: *product}
	db.Create(cartProduct)

	productCategory := &ProductCategory {Product: *product, Category: *category}
	db.Create(productCategory)

	productAdvert := &ProductAdvert {Product: *product, Advert: *advert}
	db.Create(productAdvert)

	productCollection := &ProductCollection{ Product : *product, Collection: *collection}
	db.Create(productCollection)

	// db.Commit()
	return db, nil
}

func TestFK( t *testing.T) {
	db, err := InitDb()
	if err != nil {
		t.Fatal("Could not init database")
	}
	if db == nil {
		t.Fatal("No database connection")
	}
	t.Run("testing set null of subcategory parrent after deleting parrent", func (t *testing.T) {
		category := &Category{}
		
		res := db.Where(&Category{ParrentCategory: nil}).First(category)
		if (res.Error != nil){
			t.Error(res.Error)
		}
		
		chldCategories := []Category{}
		db.Where(&Category{ParrentCategoryId: &category.ID}).Find(&chldCategories)
		
		if (len(chldCategories) == 0) {
			t.Error("No child categories")
		}
		t.Log(chldCategories)
		db.Unscoped().Delete(category)
		
		for _, curCat := range chldCategories {
			
			updCat := &Category{}
			db.First(updCat, curCat.ID)
			if (curCat.ParrentCategory != nil) {
				t.Error("Parrent category is not set to nil after deleting parrent category")
			}
			if (curCat.ParrentCategoryId != nil) {
				//t.Error("Parrent category ID is not set to nil after deleting parrent category")
			}
		}
		//t.Error()
		// db.Commit()
	})
	db, err = InitDb()
	if db == nil || err != nil {
		t.Fatal("db initialization failed")
	}
	t.Run("test cascade delete discount", func (t *testing.T) {
		discount := &Discount{}
		db.First(discount)
		// fmt.Print(discount)
		product := &Product{}
		db.First(product, discount.ProductId)
		db.Select("discounts").Delete(product)
		//db.Delete(product)
		foundDiscounts := []Discount{};
		db.Find(&foundDiscounts,&Discount{Product: *product})
		if len(foundDiscounts) != 0 {
			t.Error("Discount not deleted after product delete")
		}
	})
}
