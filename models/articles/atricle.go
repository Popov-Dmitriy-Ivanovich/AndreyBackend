package articles

import (
	"models/types"
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Html types.DbFile
	UserID uint
	ArticleMedia []ArticleMedia
}
type ArticleMedia struct {
	ID uint `gorm:"primaryKey"`
	ArticleID uint
	File types.DbFile
}

func (article *Article) AfterDelete( db *gorm.DB) error {
	res := db.Where("article_id = ?", article.ID).Delete(&ArticleMedia{})
	return res.Error
}