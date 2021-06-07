package auth

import (
	"fmt"
	//"log"
	"net/http"
	//"os"
	//"strings"
	"time"
  //"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

/*
func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Header["Token"] != nil {

      token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
          return nil, fmt.Errorf(("Invalid Signing Method"))
        }
        aud := "wordword"
        checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
        if !checkAudience {
          return nil, fmt.Errorf(("invalid aud"))
        }
        // verify iss claim
        iss := "jwtgo.io"
        checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
        if !checkIss {
          return nil, fmt.Errorf(("invalid iss"))
        }

        return MySigningKey, nil
      })
      if err != nil {
        fmt.Fprintf(w, err.Error())
      }

      if token.Valid {
        endpoint(w, r)
      }

    } else {
      fmt.Fprintf(w, "No Authorization Token provided")
    }
  })
}
*/

/////////////////////////////////------------------------------------------///////////////////////////////////////////
var MySigningKey = []byte("intel")

func IsAuthorized() gin.HandlerFunc {
  return func(c *gin.Context) {
    tokenv, err := GetJWT(c)
    if err != nil {
      c.JSON(http.StatusUnauthorized, "Error Un")
    }
    token, err := jwt.Parse(tokenv, func(token *jwt.Token) (interface{}, error) {
      return []byte(mySigningKey), nil
    })
    if err != nil {
      c.JSON(http.StatusUnauthorized, err.Error())
    }
    _, ok := token.Claims.(jwt.MapClaims)
    if ok && token.Valid {  
      claims := jwt.MapClaims{}
      token0, _ := jwt.ParseWithClaims(tokenv, claims,func(token *jwt.Token) (interface{}, error) {
        return []byte(mySigningKey), nil
      })
      fmt.Println("Hello There!\n")
      fmt.Println(claims,"\n")
      // ... error handling
      fmt.Println("---->", token0)
      y := claims["client"]
      c.Set("client", y)

    } else {
      c.Abort()
      c.JSON(http.StatusUnauthorized, "Invalid token")
    }
  }
}

/*func handleRequests() {
  http.Handle("/", isAuthorized(homePage))
  log.Fatal(http.ListenAndServe(":9001", nil))
}*/



//////////////////////////////////////////////////////////////////////////////////////////////////////
var mySigningKey = []byte("intel")

func CreateJWT(uid string, password string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
  
	claims := token.Claims.(jwt.MapClaims)
  
	claims["authorized"] = true
	claims["client"] = uid
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
  token,_ := c.Cookie("Token")
	if token == "" {
		return "", fmt.Errorf("missing tokken")
	}
	return token, nil
}
  
  



