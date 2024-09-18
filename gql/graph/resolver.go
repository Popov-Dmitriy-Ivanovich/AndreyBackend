package graph

//go:generate go run github.com/99designs/gqlgen generate
import (
	"backend/gql/graph/model"

	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	todos []*model.Todo
	db    *gorm.DB
}
