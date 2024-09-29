package main

import (
	"backend/models/products"
	"backend/models/users"
	"encoding/json"

	// "go/token"
	"graphql"
	"io"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

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
	Username string `json:"username"`
	UserId uint
	IsAdmin  bool
	jwt.RegisteredClaims
}

func generateJWT(isAdmin bool, userId uint) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: "username",
		IsAdmin:  isAdmin,
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(graphql.GetGlobalServerData().JwtKey)
}

func Login(c *gin.Context) {
	reqBody := c.Request.Body
	type LoginData struct {
		Login    string
		Password string
	}

	var body LoginData
	rowBody, errReader := io.ReadAll(reqBody)
	if errReader != nil {
		c.JSON(http.StatusTeapot, gin.H{"error": "Can't parse request body"})
		return
	}
	json.Unmarshal(rowBody, &body)

	globalServerData := graphql.GetGlobalServerData()
	user := users.User{}
	userQueryRes := globalServerData.Db.Where("login = ?", body.Login).First(&user)
	if userQueryRes.Error != nil || userQueryRes.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if body.Password != user.Password {
		c.JSON(http.StatusForbidden, gin.H{"error": "wrong password"})
		return
	}
	token, err := generateJWT(user.IsAdmin, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	globalServerData.Tokens.Mtx.Lock()
	defer globalServerData.Tokens.Mtx.Unlock()
	globalServerData.Tokens.Data = append(globalServerData.Tokens.Data, token)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		type RegisterData struct {
			Name     string
			Login    string
			Password string
			Email    string
		}
		rowBody, errReader := io.ReadAll(c.Request.Body)
		if errReader != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not parse body"})
			return
		}
		registerData := RegisterData{}
		json.Unmarshal(rowBody, &registerData)

		db := graphql.GetGlobalServerData().Db
		res := db.Create(&users.User{Name: registerData.Name, Login: registerData.Login, Password: registerData.Password, Email: registerData.Email, IsActive: true})
		if res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

func AuthMiddleware(isAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		claims := Claims{}
		tkn, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
			return graphql.GetGlobalServerData().JwtKey, nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not parse token"})
			c.Abort()
			return
		}
		if !tkn.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}
		if isAdmin && !claims.IsAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin role required"})
		}
		c.Next()
	}
}

func GetUserCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		claims := Claims{}
		tkn, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
			return graphql.GetGlobalServerData().JwtKey, nil
		})
		if err != nil || !tkn.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error":"unauth"})
			return
		}
		userId := claims.UserId
		db := graphql.GetGlobalServerData().Db
		cart := products.Cart{}
		cartRes := db.Where("user_id = ?", userId).First(&cart)
		if cartRes.Error != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"could not find user cart"})
			return
		}
		prods := []*products.Product{}
		prodRes := db.Model(&products.Product{}).Joins("inner join product_carts on products.id == product_carts.product_id").Where("cart_id = ?", cart.ID).Find(&prods)
		if prodRes.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"db"})
		}
		c.JSON(http.StatusOK, prods)
	}
}

func AddProductToCart() gin.HandlerFunc {
	return func (c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		claims := Claims{}
		tkn, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
			return graphql.GetGlobalServerData().JwtKey, nil
		})
		if err != nil || !tkn.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error":"unauth"})
			return
		}
		userId := claims.UserId
		db := graphql.GetGlobalServerData().Db
		
		cart := products.Cart{}
		res := db.Where("user_id = ?", userId).First(&cart)
		if res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H {"error":"cart not found"})
			return
		}
		type bodyScheme struct {
			ProductId uint
		}
		ginBody := c.Request.Body;
		rowBody, errReader := io.ReadAll(ginBody)
		if errReader != nil {
			c.JSON(http.StatusTeapot, gin.H{"error": "Can't parse request body"})
			return
		}
		body := bodyScheme{}
		json.Unmarshal(rowBody, &body)
		prodcutToAdd := products.Product{}
		prodRes := db.First(&prodcutToAdd,body.ProductId)
		if prodRes.Error != nil || prodRes.RowsAffected == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error"})
			return
		}
		cart.Product = append(cart.Product, &prodcutToAdd)
		saveRes := db.Save(&cart)
		if saveRes.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error saving cart"})
			return
		}
		c.JSON(http.StatusOK,gin.H{})
	}
}

func PurcacheUserCart() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func CreateReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		type RequestBody struct {
			ProductID uint
			UserID    uint
			Text      string
			Stars     uint
		}
		ginBody := c.Request.Body
		rowBody, errReader := io.ReadAll(ginBody)
		if errReader != nil {
			c.JSON(http.StatusTeapot, gin.H{"error": "Can't parse request body"})
			return
		}
		body := RequestBody{}
		json.Unmarshal(rowBody, &body)
		user := users.User{}
		db := graphql.GetGlobalServerData().Db
		res := db.First(&user, body.UserID)

		if res.Error != nil || res.RowsAffected == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		user.Reviews = append(user.Reviews, products.Review{
			ProductID: body.ProductID,
			UserID: user.ID,
			Text: body.Text,
			Stars: body.Stars,
		})
		resSave := db.Save(&user)
		if resSave.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error" : "could not save user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

func GetReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()
		productId := query.Get("productId")
		db := graphql.GetGlobalServerData().Db
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
	// globalServerData := graphql.GetGlobalServerData()
	resolver := gin.Default()

	adminResolver := resolver.Group("/Admin", AuthMiddleware(true))
	userResolver := resolver.Group("/User", AuthMiddleware(false))
	unauthResolver := resolver.Group("/")

	userResolver.POST("/UserQuery", graphqlUserHandler())
	userResolver.POST("/CreateReview", CreateReview())
	userResolver.GET("/GetCart", GetUserCart())
	userResolver.POST("/AddProduct", AddProductToCart())
	adminResolver.POST("/AdminQuery", graphqlAdminHandler())

	unauthResolver.GET("/GQLPlayground", graphqlPlaygroundHandler())
	unauthResolver.POST("/Login", Login)
	unauthResolver.POST("/Register", Register())
	unauthResolver.GET("/Review", GetReview())
	resolver.Run()
}
