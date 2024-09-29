package users

import (
	"backend/models/articles"
	"backend/models/products"

	"gorm.io/gorm"
)

type User struct {

	gorm.Model
	Name     string
	Login    string `gorm:"unqique"`
	Email    string `binding:"required,email"`
	IsActive bool
	IsAdmin  bool
	Password string
	Articles []articles.Article
	Reviews  []products.Review
	Cart 	 *products.Cart
}

func (user *User) AfterCreate (db *gorm.DB) error {
	cart := &products.Cart{}
	user.Cart = cart
	res := db.Save(user)
	return res.Error
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
