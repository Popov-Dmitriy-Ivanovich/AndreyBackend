package articles

import (
	// "models/products"
	"models/types"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestProductCategoryDelete(t *testing.T) { 
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil || db == nil {
		//panic("failed to connect database")
		t.Fatal("DB INIT FAILED")
	}

	db.AutoMigrate(&Article{}, &ArticleMedia{})
	articleMedia := []ArticleMedia{ArticleMedia{File: types.DbFile{}}}
	article := Article{ArticleMedia: articleMedia}
	db.Create(&article)
	res := db.Find(&[]ArticleMedia{})
	if res.RowsAffected == 0 {
		t.Error("No article media")
	}
	db.Delete(&article)
	res = db.Find(&[]ArticleMedia{})
	if res.RowsAffected != 0 {
		t.Error("Article media has not deleted after deleting article")
	}
}