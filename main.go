package main

import (
	"backend/models/products"
	"backend/models/users"
	"encoding/json"
	"errors"

	// "go/token"
	"graphql"
	"io"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

	"time"
)

// "backend/db"
// "backend/gql"

func graphqlUserHandler() gin.HandlerFunc {
	hndl := graphql.GetGlobalServerData().GqlUser
	return func(c *gin.Context) {
		hndl.ServeHTTP(c.Writer, c.Request)
	}
}

func graphqlAdminHandler() gin.HandlerFunc {
	hndl := graphql.GetGlobalServerData().GqlAdmin
	return func(c *gin.Context) {
		hndl.ServeHTTP(c.Writer, c.Request)
	}
}

func graphqlPlaygroundHandler() gin.HandlerFunc {
	pg := playground.Handler("GraphQL", "/UserQuery")
	return func(c *gin.Context) {
		pg.ServeHTTP(c.Writer, c.Request)
	}
}

type Claims struct {
	// Username string `json:"username"`
	UserId  uint
	IsAdmin bool
	jwt.RegisteredClaims
}

type ServerRoutes struct {
	globalServerData *graphql.GlobalSingletonServerData
	db               *gorm.DB
	jwtKey           *[]byte
}

func (sr *ServerRoutes) getBody(c *gin.Context, result any) error {
	//ginBody := c.Request.Body
	rowBody, errReader := io.ReadAll(c.Request.Body)
	if errReader != nil {
		return errReader
	}
	//errUnmarshal := 
	
	return json.Unmarshal(rowBody, result)
}

func (sr *ServerRoutes) getClaims(c *gin.Context) (Claims, error) {
	token := c.Request.Header.Get("Authorization")
	claims := Claims{}
	tkn, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return sr.jwtKey, nil
	})
	if err != nil {
		return claims, err
	}
	if !tkn.Valid {
		return claims, errors.New("auth token is invalid")
	}
	return claims, nil
}

func CreateServerRoutes() (ServerRoutes, error) {
	globData := graphql.GetGlobalServerData()
	sr := ServerRoutes{
		globalServerData: globData,
		db:               globData.Db,
		jwtKey:           &globData.JwtKey,
	}
	return sr, nil
}

func (sr *ServerRoutes) generateJWT(isAdmin bool, userId uint) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		IsAdmin: isAdmin,
		UserId:  userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(sr.jwtKey)
}

func (sr *ServerRoutes) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		type LoginData struct {
			Login    string
			Password string
		}
		var body LoginData
		errGettingBody := sr.getBody(c, &body)
		if errGettingBody != nil {
			c.JSON(http.StatusInternalServerError, errGettingBody)
			return
		}

		user := users.User{}
		userQueryRes := sr.db.Where("login = ?", body.Login).First(&user)
		if userQueryRes.Error != nil || userQueryRes.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if body.Password != user.Password {
			c.JSON(http.StatusForbidden, gin.H{"error": "wrong password"})
			return
		}
		token, err := sr.generateJWT(user.IsAdmin, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}
		sr.globalServerData.Tokens.Mtx.Lock()
		defer sr.globalServerData.Tokens.Mtx.Unlock()
		sr.globalServerData.Tokens.Data = append(sr.globalServerData.Tokens.Data, token)
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func (sr *ServerRoutes) 	Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		type RegisterData struct {
			Name     string
			Login    string
			Password string
			Email    string
		}

		registerData := RegisterData{}
		err := sr.getBody(c, &registerData)
		if err != nil {
			c.JSON(http.StatusInternalServerError,err)
			return
		}
		db := sr.db
		res := db.Create(&users.User{Name: registerData.Name, Login: registerData.Login, Password: registerData.Password, Email: registerData.Email, IsActive: true})
		if res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

func (sr *ServerRoutes) AuthMiddleware(isAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := sr.getClaims(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}
		if isAdmin && !claims.IsAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin role required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (sr *ServerRoutes) GetUserCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := sr.getClaims(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err)
			return
		}

		userId := claims.UserId
		db := sr.db
		cart := products.Cart{}
		cartRes := db.Where("user_id = ?", userId).First(&cart)
		if cartRes.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not find user cart"})
			return
		}
		prods := []*products.Product{}
		prodRes := db.Model(&products.Product{}).Joins("inner join product_carts on products.id == product_carts.product_id").Where("cart_id = ?", cart.ID).Find(&prods)
		if prodRes.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db"})
		}
		c.JSON(http.StatusOK, prods)
	}
}

func (sr *ServerRoutes) AddProductToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := sr.getClaims(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err)
		}
		userId := claims.UserId
		db := sr.db

		cart := products.Cart{}
		res := db.Where("user_id = ?", userId).First(&cart)
		if res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cart not found"})
			return
		}
		type bodyScheme struct {
			ProductId uint
		}
		
		body := bodyScheme{}
		sr.getBody(c, body)

		prodcutToAdd := products.Product{}
		prodRes := db.First(&prodcutToAdd, body.ProductId)
		if prodRes.Error != nil || prodRes.RowsAffected == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error"})
			return
		}
		cart.Product = append(cart.Product, &prodcutToAdd)
		saveRes := db.Save(&cart)
		if saveRes.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving cart"})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}

func PurcacheUserCart() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (sr *ServerRoutes) CreateReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		type RequestBody struct {
			ProductID uint
			UserID    uint
			Text      string
			Stars     uint
		}
		body := RequestBody{}
		err := sr.getBody(c, &body)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err)
		}

		user := users.User{}
		db := sr.db
		res := db.First(&user, body.UserID)

		if res.Error != nil || res.RowsAffected == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		user.Reviews = append(user.Reviews, products.Review{
			ProductID: body.ProductID,
			UserID:    user.ID,
			Text:      body.Text,
			Stars:     body.Stars,
		})
		resSave := db.Save(&user)
		if resSave.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

func (sr *ServerRoutes) GetReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()
		productId := query.Get("productId")
		db := sr.db
		review := []products.Review{}
		res := db.Where("product_id = ?", productId).Find(&review)
		if res.Error != nil {
			c.JSON(http.StatusInternalServerError, "db error")
			return
		}
		c.JSON(http.StatusOK, review)
	}
}

func main() {
	graphql.InitGlobalServerData()
	
	sr, error := CreateServerRoutes()

	if error != nil {
		panic("could not create server routes object")
	}

	resolver := gin.Default()

	adminResolver := resolver.Group("/Admin", sr.AuthMiddleware(true))
	userResolver := resolver.Group("/User", sr.AuthMiddleware(false))
	unauthResolver := resolver.Group("/")

	unauthResolver.POST("/UserQuery", graphqlUserHandler())
	userResolver.POST("/CreateReview", sr.CreateReview())
	userResolver.GET("/GetCart", sr.GetUserCart())
	userResolver.POST("/AddProduct", sr.AddProductToCart())
	adminResolver.POST("/AdminQuery", graphqlAdminHandler())

	unauthResolver.GET("/GQLPlayground", graphqlPlaygroundHandler())
	unauthResolver.POST("/Login", sr.Login())
	unauthResolver.POST("/Register", sr.Register())
	unauthResolver.GET("/Review", sr.GetReview())
	resolver.Run()
}
