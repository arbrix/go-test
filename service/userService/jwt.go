package userService

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/arbrix/go-test/config"
	"github.com/arbrix/go-test/model"
	"github.com/arbrix/go-test/util/log"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreateToken(user *model.User) (string, int, error) {
	// Create the token
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	// Set some claims
	token.Claims["AccessToken"] = "level1"
	token.Claims["CustomUserInfo"] = struct {
		Name string
		Kind string
	}{user.Name, "human"}
	token.Claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(config.SecretKey)
	if err != nil {
		log.Fatalf("Token Signing error: %v\n", err)
		return "", http.StatusInternalServerError, errors.New("Sorry, error while Signing Token!")
	}
	return tokenString, http.StatusOK, nil
}

func CheckToken(c *gin.Context) (bool, error) {
	authToken, err := fromAuthHeader(c)
	if err != nil {
		return false, err
	}
	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
	return token.Valid, err
}

// FromAuthHeader is a "TokenExtractor" that takes a give request and extracts
// the JWT token from the Authorization header.
func fromAuthHeader(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}
