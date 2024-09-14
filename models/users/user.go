package users

import (
	"models/articles"
	"models/products"

	"gorm.io/gorm"
)

type User struct {

	gorm.Model
	Name     string
	Login    string
	Email    string `binding:"required,email"`
	IsActive bool
	IsAdmin  bool
	Password string
	Articles []articles.Article
	Reviews  []products.Review
	Cart 	 *products.Cart
}

func (user *User) AfterDelete(db *gorm.DB) error {
	res := db.Model(&articles.Article{}).Where("user_id = ?", user.ID).Update("user_id", nil)
	if res.Error != nil {
		return res.Error
	}

	res = db.Model(&products.Cart{}).Where("user_id = ?", user.ID).Delete(&products.Cart{})
	if res.Error != nil {
		return res.Error
	}

	res = db.Model(&products.Review{}).Where("user_id = ?", user.ID).Delete(&products.Review{})
	
	return res.Error
}
