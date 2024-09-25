package main

import (
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
	IsAdmin bool
	jwt.RegisteredClaims
}

func generateJWT(isAdmin bool) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: "username",
		IsAdmin: isAdmin,
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
	token, err := generateJWT(user.IsAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"error": "Could not generate token"})
		return
	}
	globalServerData.Tokens.Mtx.Lock()
	defer globalServerData.Tokens.Mtx.Unlock()
	globalServerData.Tokens.Data = append(globalServerData.Tokens.Data, token)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Register() gin.HandlerFunc {
	return func (c *gin.Context) {
		type RegisterData struct {
			Name string
			Login    string
			Password string
			Email 	 string
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
	return func (c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		claims := Claims{}
		tkn, err := jwt.ParseWithClaims(token,&claims, func(token *jwt.Token) (interface{}, error) {
            return graphql.GetGlobalServerData().JwtKey, nil
        })
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Could not parse token"})
			return
		}
		if !tkn.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
                "error": "unauthorized",
            })
			return
		}
		if isAdmin && !claims.IsAdmin {
			c.JSON(http.StatusUnauthorized, gin.H {"error":"Admin role required"})
		}
		c.Next()
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
	adminResolver.POST("/AdminQuery", graphqlAdminHandler())
	
	unauthResolver.GET("/GQLPlayground", graphqlPlaygroundHandler())
	unauthResolver.POST("/Login", Login)
	unauthResolver.POST("/Register", Register())
	resolver.Run()
}
