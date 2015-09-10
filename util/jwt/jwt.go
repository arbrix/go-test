package jwt

import (
	"errors"
	"github.com/arbrix/go-test/app"
	"github.com/arbrix/go-test/model"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strings"
	"time"
)

type Token struct {
}

//Create is make new jwt and set it to context
func (t *Token) Create(ec *echo.Context, a *app.App, usr *model.User) (int, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["id"] = usr.ID
	token.Claims["email"] = usr.Email
	token.Claims["name"] = usr.Name
	token.Claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	scrKey, err := t.getSecretKey(a)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(scrKey))
	if err != nil {
		log.Fatalf("Token Signing error: %v\n", err)
		return http.StatusInternalServerError, errors.New("Sorry, error while Signing Token!")
	}
	ec.Set("jwt", tokenString)
	return http.StatusOK, nil
}

// Parse is a "TokenExtractor" & "TokenParser" that takes a give request and extracts
// the JWT token from the Authorization header and parse it.
func (t *Token) Parse(ec *echo.Context, a *app.App) (*jwt.Token, error) {
	authHeader := ec.Request().Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("JWT is not in Header, it is empty")
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		log.Printf("Authorization header format must be Bearer {token}")
		return nil, errors.New("JWT is not in Header")
	}
	scrKey, err := t.getSecretKey(a)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(authHeaderParts[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(scrKey), nil
	})
	return token, nil
}

func (t *Token) getSecretKey(a *app.App) (string, error) {
	scrKey, err := a.GetConfig().Get("SecretKey")
	if err != nil {
		log.Fatalf("Secret not defined in config: %v\n", err)
		return "", errors.New("Sorry, important options is not set via config files, please check it!")
	}
	return scrKey.(string), nil
}
