package productmedia

import (
	"models/types"

	"gorm.io/gorm"
)

type ProductMedia struct {
	gorm.Model
	ProductID uint
	File types.DbFile
}
