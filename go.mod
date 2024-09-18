module backend

go 1.22.5

replace backend/db => ./db

replace backend/gql => ./gql

require backend/gql v0.0.0-00010101000000-000000000000

require (
	github.com/99designs/gqlgen v0.17.51 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/vektah/gqlparser/v2 v2.5.16 // indirect
	golang.org/x/text v0.18.0 // indirect
	gorm.io/gorm v1.25.12 // indirect
)
