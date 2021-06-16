package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

/////////////////////////////////------------------------------------------///////////////////////////////////////////
var MySigningKey = []byte(os.Getenv("MYCODE"))

func IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenv, err := GetJWT(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "Error UnAuthorized")
			c.Abort()
			return
		}
		token, er := jwt.Parse(tokenv, func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		})
		if er != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
		}
		_, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			claims := jwt.MapClaims{}
			jwt.ParseWithClaims(tokenv, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(mySigningKey), nil
			})

			y := claims["client"]
			c.Set("client", y)

		} else {
			c.Abort()
			c.JSON(http.StatusUnauthorized, "Invalid token")
		}
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
var mySigningKey = []byte(os.Getenv("MYCODE"))

func CreateJWT(IDU string, password string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = IDU
	claims["aud"] = password
	claims["iss"] = "jwtgo.io"
	claims["exp"] = time.Now().Add(time.Minute * 100).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", fmt.Errorf("something went wrong")
	}

	return tokenString, nil
}

func GetJWT(c *gin.Context) (string, error) {
	token, err := c.Cookie("Token")
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to access token"})
	}
	if token == "" {
		return "", fmt.Errorf("missing tokken")
	}
	return token, nil
}
