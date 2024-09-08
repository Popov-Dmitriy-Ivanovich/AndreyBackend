package models
import (
	"gorm.io/gorm"
	"time"
)

type Discount struct {
	gorm.Model
	ProductID      *uint
	NewPrice       PriceType
	Style          string
	ExpirationDate time.Time
	// Product        Product `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}