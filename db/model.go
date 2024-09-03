package db

import (
	"database/sql/driver"
	"errors"

	"log"
	"os"
	"time"

	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type PriceType float32

type DbFile struct {
	Path string
}

func (dbFile DbFile) AsFile() (*os.File, error) {
	if dbFile.Path == "" {
		return nil, nil
	}
	file, err := os.Open(dbFile.Path) // For read access.
	if err != nil {
		log.Fatal(err)
	}
	return file, err
}

func (dbFile *DbFile) Scan(value any) error {

	str, ok := value.(string)
	if !ok {
		log.Fatal(ok)
		return errors.New("failed to unmarshal file value")
	}
	dbFile.Path = str
	return nil
}
func (dbFile DbFile) Value() (driver.Value, error) {
	_, err := dbFile.AsFile()
	if err != nil {
		return nil, errors.New("could not open file")
	}
	return dbFile.Path, nil
}

type Product struct {
	gorm.Model
	Name        string
	Desctiption string
	Price       PriceType
	IsActive    bool
	Count       uint64
	Picture     DbFile
}

type Discount struct {
	gorm.Model
	ProductId      int
	NewPrice       PriceType
	Style          string
	ExpirationDate time.Time
	Product        Product `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type ProductMedia struct {
	gorm.Model
	ProductId int
	Product   Product `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	File      DbFile
}

type Category struct {
	gorm.Model
	Name              string
	Description       string
	ParrentCategoryId *uint
	ParrentCategory   *Category `gorm:"foreignKey:ParrentCategoryId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Picture           DbFile
}

type Advert struct {
	gorm.Model
	ExpirationDate time.Time
	Text           string
	Style          string
}

type User struct {
	gorm.Model
	Name     string
	Login    string
	Email    string `binding:"required,email"`
	IsActive bool
	IsAdmin  bool
	Password string
}

type Article struct {
	gorm.Model
	Html     DbFile
	AuthorId int
	Author   User `gorm:"foreignKey=AuthorId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type ArticleMedia struct {
	gorm.Model
	ArticleId int
	Article   Article `gorm:"foreignKey=ArticleId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	File      DbFile
}

type Collection struct {
	gorm.Model
	Name        string
	Description string
	Picture     DbFile
}

type Cart struct {
	gorm.Model
	UserId int
	User   User `gorm:"foreignKey=UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type NewsSubscribers struct {
	gorm.Model
	Email string `binding:"required,email"`
}

type ProductReview struct {
	gorm.Model
	ProductId int
	Product   Product `gorm:"foreignKey=ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserId    int
	User      User `gorm:"foreignKey=UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Text      string
	Stars     uint8
}

type CartProduct struct {
	gorm.Model
	CartId    uint
	Cart      Cart `gorm:"foreignKey=CartId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ProductId uint
	Product   Product `gorm:"foreignKey=ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type ProductCategory struct {
	gorm.Model
	ProductId  uint
	Product    Product `gorm:"foreignKey=ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CategoryId uint
	Category   Category `gorm:"foreignKey=CategoryId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type ProductAdvert struct {
	gorm.Model
	ProductId uint
	Product   Product `gorm:"foreignKey=ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AdvertId  uint
	Advert    Advert `gorm:"foreignKey=AdvertId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Priority  uint
}

type ProductCollection struct {
	gorm.Model
	ProductId    uint
	Product      Product `gorm:"foreignKey=ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CollectionId uint
	Collection   Collection `gorm:"foreignKey=CollectionId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
