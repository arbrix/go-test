package jwt

import (
	"errors"
	"github.com/arbrix/go-test/interfaces"
	"github.com/arbrix/go-test/model"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strings"
	"time"
)

type Tokenizer struct {
	a interfaces.App
}

func NewTokenizer(a interfaces.App) *Tokenizer {
	t := &Tokenizer{a: a}
	return t
}

//Create is make new jwt and set it to context
func (t *Tokenizer) Create(ec *echo.Context, usr *model.User) (int, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["id"] = usr.ID
	token.Claims["email"] = usr.Email
	token.Claims["name"] = usr.Name
	token.Claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	scrKey, err := t.getSecretKey()
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
func (t *Tokenizer) Parse(ec *echo.Context) (*jwt.Token, error) {
	authHeader := ec.Request().Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("JWT is not in Header, it is empty")
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		log.Printf("Authorization header format must be Bearer {token}")
		return nil, errors.New("JWT is not in Header")
	}
	scrKey, err := t.getSecretKey()
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(authHeaderParts[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(scrKey), nil
	})
	return token, nil
}

func (t *Tokenizer) Check() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ec *echo.Context) error {
			token, err := t.Parse(ec)
			if err != nil {
				log.Println(err.Error())
				ec.Error(echo.NewHTTPError(http.StatusInternalServerError))
			}
			if token.Valid == false {
				ec.Error(echo.NewHTTPError(http.StatusForbidden))
				return nil
			}
			if err := h(ec); err != nil {
				ec.Error(err)
			}
			return nil
		}
	}
}

func (t *Tokenizer) getSecretKey() (string, error) {
	scrKey, err := t.a.GetConfig().Get("secret")
	if err != nil {
		log.Fatalf("Secret not defined in config: %v\n", err)
		return "", errors.New("Sorry, important options is not set via config files, please check it!")
	}
	return scrKey.(string), nil
}
