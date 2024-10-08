package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.53

import (
	"backend/models/discounts"
	"backend/models/productmedia"
	"backend/models/products"
	"backend/models/types"
	"context"
	"errors"
	"fmt"
	"graphql/graph/model"
	"strconv"
	"time"
)

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented: CreateTodo - createTodo"))
}

// CreateDiscount is the resolver for the createDiscount field.
func (r *mutationResolver) CreateDiscount(ctx context.Context, input model.NewDiscount) (*model.Discount, error) {
	if r.IsAdmin != true {
		return nil, errors.New("access denied")
	}
	prodId, err := strconv.ParseUint(input.ProductID, 10, 64)
	if err != nil {
		return nil, err
	}

	prod := products.Product{}
	res := r.Db.First(&prod, prodId)
	if res.Error != nil {
		return nil, res.Error
	}

	newDiscont := discounts.Discount{NewPrice: types.PriceType(input.NewPrice), Style: input.Style}
	prod.Discounts = append(prod.Discounts, newDiscont)
	updateRes := r.Db.Save(&prod)
	return &model.Discount{ID: "value"}, updateRes.Error
}

// CreateProduct is the resolver for the createProduct field.
func (r *mutationResolver) CreateProduct(ctx context.Context, input *model.NewProduct) (*model.Product, error) {
	if r.IsAdmin != true {
		return nil, errors.New("access denied")
	}
	foundCat := []*products.Category{}
	r.Db.Where("id in ?", input.CategoriesID).Find(&foundCat)

	newProd := products.Product{
		Name:        input.Name,
		Desctiption: input.Description,
		Price:       types.PriceType(input.Count),
		IsActive:    input.IsActive,
		Count:       uint64(input.Count),
		Picture:     types.DbFile{Path: input.Picture},
		Categories:  foundCat,
	}
	db := r.Db
	res := db.Create(&newProd)
	return &model.Product{ID: fmt.Sprintf("%d", newProd.ID)}, res.Error
}

// CreateCollection is the resolver for the createCollection field.
func (r *mutationResolver) CreateCollection(ctx context.Context, input *model.NewCollection) (*model.Collection, error) {
	if r.IsAdmin != true {
		return nil, errors.New("access denied")
	}
	prodIds := make([][]uint, 0, len(input.ProductIds))
	for _, el := range input.ProductIds {
		if el == nil {
			return nil, errors.New("undefinded product id")
		}
		id, err := strconv.ParseUint(*el, 10, 64)
		if err != nil {
			return nil, errors.New("could not parse product id")
		}
		prodIds = append(prodIds, []uint{uint(id)})
	}
	foundProd := []*products.Product{}
	res := r.Db.Where("id IN ?", prodIds).Find(&foundProd)
	if res.Error != nil {
		return nil, res.Error
	}
	newCollection := products.Collection{
		Name:        input.Name,
		Description: input.Description,
		Picture:     types.DbFile{Path: input.Picture},
		Products:    foundProd}
	r.Db.Create(&newCollection)
	return &model.Collection{ID: fmt.Sprintf("%d", newCollection.ID)}, nil
}

// CreateProductMedia is the resolver for the createProductMedia field.
func (r *mutationResolver) CreateProductMedia(ctx context.Context, input *model.NewProductMedia) (*model.ProductMedia, error) {
	if r.IsAdmin != true {
		return nil, errors.New("access denied")
	}
	id, err := strconv.ParseUint(input.ProductID, 10, 64)
	if err != nil {
		return nil, err
	}
	foundProduct := products.Product{}
	res := r.Db.Find(&foundProduct, id)

	if res.Error != nil {
		return nil, res.Error
	}
	prodMedia := productmedia.ProductMedia{File: types.DbFile{Path: input.File}}
	foundProduct.Media = append(foundProduct.Media, prodMedia)
	res = r.Db.Save(&foundProduct)
	if res.Error != nil {
		return nil, res.Error
	}
	return &model.ProductMedia{Product: &model.Product{ID: fmt.Sprintf("%d", foundProduct.ID)}, File: input.File}, nil
}

// CreateAdvert is the resolver for the createAdvert field.
func (r *mutationResolver) CreateAdvert(ctx context.Context, input *model.NewAdvert) (*model.Advert, error) {
	if r.IsAdmin != true {
		return nil, errors.New("access denied")
	}
	foundProd := []*products.Product{}
	res := r.Db.Where("id in ?", input.ProductIds).Find(&foundProd)
	if res.Error != nil {
		return nil, res.Error
	}
	expDate, err := time.Parse("2006-01-02", input.ExpirationDate)
	if err != nil {
		return nil, err
	}
	advert := products.Advert{
		ExpirationDate: expDate,
		AdvertText:     input.Text,
		Style:          input.Style,
		Products:       foundProd,
	}
	r.Db.Create(&advert)
	return &model.Advert{Style: "test"}, nil
}

// CreateCategory is the resolver for the createCategory field.
func (r *mutationResolver) CreateCategory(ctx context.Context, input *model.NewCategory) (*model.Category, error) {
	if r.IsAdmin != true {
		return nil, errors.New("access denied")
	}
	var parentId *uint
	parentId = nil
	if input.ParentID != nil {
		convId, err := strconv.ParseUint(*input.ParentID, 10, 64)
		if err != nil {
			return nil, err
		}
		resId := uint(convId)
		parentId = &resId
	}
	newCat := products.Category{
		Name:        input.Name,
		Description: input.Description,
		Picture:     types.DbFile{Path: input.Picture},
		ParentID:    parentId,
	}
	r.Db.Create(&newCat)
	return &model.Category{ID: fmt.Sprintf("%d", newCat.ID)}, nil
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented: Todos - todos"))
}

// Discounts is the resolver for the discounts field.
func (r *queryResolver) Discounts(ctx context.Context) ([]*model.Discount, error) {
	db := r.Db
	foundDiscounts := []discounts.Discount{}
	res := db.Find(&foundDiscounts)
	if res.Error != nil {
		return nil, res.Error
	}
	gqlDiscounts := make([]*model.Discount, 0, len(foundDiscounts))
	for _, el := range foundDiscounts {
		gqlDiscounts = append(gqlDiscounts, &model.Discount{
			ID:             fmt.Sprintf("%d", el.ID),
			ProductID:      fmt.Sprintf("%d", *el.ProductID),
			NewPrice:       float64(el.NewPrice),
			Style:          el.Style,
			ExpirationDate: el.ExpirationDate.String(),
		})
	}
	return gqlDiscounts, nil
}

// Products is the resolver for the products field.
func (r *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {
	db := r.Db
	foundProducts := []products.Product{}
	res := db.Find(&foundProducts)
	if res.Error != nil {
		return nil, res.Error
	}
	returnedProducts := make([]*model.Product, 0, len(foundProducts))
	for _, el := range foundProducts {
		returnedProducts = append(returnedProducts,
			&model.Product{
				ID:          fmt.Sprintf("%d", el.ID),
				Name:        el.Name,
				Description: el.Desctiption,
				Price:       float64(el.Price),
				IsActive:    el.IsActive,
				Count:       int(el.Count),
				Picture:     el.Picture.Path})
	}
	return returnedProducts, res.Error
}

// Collections is the resolver for the collections field.
func (r *queryResolver) Collections(ctx context.Context) ([]*model.Collection, error) {
	foundCollections := []products.Collection{}
	res := r.Db.Find(&foundCollections)
	if res.Error != nil {
		return nil, res.Error
	}
	resCollections := make([]*model.Collection, 0, len(foundCollections))
	for _, el := range foundCollections {
		resCollections = append(resCollections, &model.Collection{
			ID:          fmt.Sprintf("%d", el.ID),
			Name:        el.Name,
			Description: el.Description,
			Picture:     el.Picture.Path,
		})
	}
	return resCollections, nil
}

// ProductMedias is the resolver for the productMedias field.
func (r *queryResolver) ProductMedias(ctx context.Context) ([]*model.ProductMedia, error) {
	foundMedia := []productmedia.ProductMedia{}
	res := r.Db.Find(&foundMedia)
	if res.Error != nil {
		return nil, res.Error
	}
	returnedMedia := make([]*model.ProductMedia, 0, len(foundMedia))
	for _, el := range foundMedia {
		prod := &products.Product{}
		res = r.Db.First(prod, el.ID)
		if res.Error != nil {
			returnedMedia = append(returnedMedia, &model.ProductMedia{
				File:    el.File.Path,
				Product: nil,
			})
			continue
		}
		returnedMedia = append(returnedMedia, &model.ProductMedia{
			File:    el.File.Path,
			Product: &model.Product{ID: fmt.Sprintf("%d", prod.ID)},
		})
	}
	return returnedMedia, nil
}

// Adverts is the resolver for the adverts field.
func (r *queryResolver) Adverts(ctx context.Context) ([]*model.Advert, error) {
	adverts := []products.Advert{}
	res := r.Db.Find(&adverts)
	if res.Error != nil {
		return nil, res.Error
	}
	returnedAdverts := make([]*model.Advert, 0, len(adverts))
	for _, el := range adverts {
		foundProd := []products.Product{}
		res := r.Db.Table("products").Joins("inner join advert_products on advert_products.product_id = products.id").Where("advert_id = ?", el.ID).Find(&foundProd)
		if res.Error != nil {
			returnedAdverts = append(returnedAdverts, &model.Advert{
				Text:  el.AdvertText,
				Style: el.Style,
				//Products: modelsProd,
			})
			continue
		}

		modelsProd := make([]*model.Product, 0, len(foundProd))
		for _, el := range foundProd {
			modelsProd = append(modelsProd, &model.Product{
				ID:          fmt.Sprintf("%d", el.ID),
				Name:        el.Name,
				Description: el.Desctiption,
				Price:       float64(el.Price),
				IsActive:    el.IsActive,
				Count:       int(el.Count),
				Picture:     el.Picture.Path,
			})
		}
		returnedAdverts = append(returnedAdverts, &model.Advert{
			Text:     el.AdvertText,
			Style:    el.Style,
			Products: modelsProd,
		})
	}
	return returnedAdverts, nil
}

// Category is the resolver for the category field.
func (r *queryResolver) Category(ctx context.Context) ([]*model.Category, error) {
	foundCategories := []products.Category{}
	res := r.Db.Find(&foundCategories)

	if res.Error != nil {
		return nil, res.Error
	}

	resCategory := make([]*model.Category, 0, len(foundCategories))
	for _, el := range foundCategories {
		foundProducts := []products.Product{}

		res := r.Db.Table("products").Joins("inner join product_categories on products.id = product_categories.product_id").Where("category_id = ?", el.ID).Find(&foundProducts)

		if res.Error != nil {
			return nil, res.Error
		}

		resProducts := make([]*model.Product, 0, len(foundProducts))
		for _, prod := range foundProducts {
			resProducts = append(resProducts, &model.Product{
				ID:          fmt.Sprintf("%d", prod.ID),
				Name:        prod.Name,
				Description: prod.Desctiption,
				Price:       float64(prod.Price),
				Count:       int(prod.Count),
				IsActive:    prod.IsActive,
				Picture:     prod.Picture.Path,
			})
		}
		parentID := "null"
		if el.ParentID != nil {
			parentID = fmt.Sprintf("%d", *el.ParentID)
		}
		resCategory = append(resCategory, &model.Category{
			ID:          fmt.Sprintf("%d", el.ID),
			Name:        el.Name,
			Description: el.Description,
			ParentID:    parentID,
			Picture:     el.Picture.Path,
			Products:    resProducts,
		})
	}
	// /resCategory := []*model.Category{}
	return resCategory, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
