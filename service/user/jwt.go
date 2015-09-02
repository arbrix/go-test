package user

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/arbrix/go-test/model"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var JwtToken string

func CreateToken(userContact string) (int, error) {
	// Create the token
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	// Set some claims
	token.Claims["AccessToken"] = "level1"
	token.Claims["CustomUserInfo"] = userContact
	token.Claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		log.Fatalf("Token Signing error: %v\n", err)
		return http.StatusInternalServerError, errors.New("Sorry, error while Signing Token!")
	}
	JwtToken = tokenString
	return http.StatusOK, nil
}

func GetToken() (string, error) {
	if JwtToken == "" {
		return "", errors.New("JWT not created!")
	}
	return JwtToken, nil
}

func DecodeToken(encToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
	return token, err
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

// CurrentUser get a current user.
func CurrentUser(c *gin.Context) (model.User, error) {
	var user model.User
	encToken, err := fromAuthHeader(c)
	if err != nil {
		return user, err
	}
	jwtToken, err := DecodeToken(encToken)
	if err != nil {
		return user, err
	}
	if jwtToken.Valid == false {
		return user, errors.New("JWT not valid")
	}

	if db.ORM.Select(config.UserPublicFields+", email").Where("token = ?", jwtToken.Claims["CustomUserInfo"]).First(&user).RecordNotFound() {
		return user, errors.New("User is not found by jwt")
	}
	return user, nil
}
