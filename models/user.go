package models
import "gorm.io/gorm"

type User struct {

	gorm.Model
	Name     string
	Login    string
	Email    string `binding:"required,email"`
	IsActive bool
	IsAdmin  bool
	Password string

}