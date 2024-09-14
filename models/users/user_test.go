package users

import (
	// "models/products"
	"models/types"
	"testing"
	"models/articles"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDeletingUserArticles(t *testing.T) { 
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil || db == nil {
		//panic("failed to connect database")
		t.Fatal("DB INIT FAILED")
	}

	db.AutoMigrate(&User{}, &articles.Article{})
	article := articles.Article{Html: types.DbFile{Path: ""}}
	user := User{Name: "test user", Articles: []articles.Article{article}}
	db.Create(&user)
	
	foundArticles := []articles.Article{}
	db.Where("user_id = ?", user.ID).Find(&foundArticles)
	if len(foundArticles) == 0 {
		t.Error("Article for user not found")
	}
	
	db.Delete(&user)

	db.Where("user_id = ?", user.ID).Find(&foundArticles)
	if len(foundArticles) != 0 {
		t.Error("Article for deleted user found")
	}
}