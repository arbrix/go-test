package jwt

import (
	"github.com/arbrix/go-test/app"
	"github.com/arbrix/go-test/model"
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	user      = model.User{ID: 1, Email: "test@user.com", Name: "testUser"}
	a         = app.NewApp(&app.TestConfig{}, &app.TestOrm{})
	tokenizer = Token{}
	req, _    = http.NewRequest(echo.GET, "/", nil)
	rec       = httptest.NewRecorder()
	c         = echo.NewContext(req, echo.NewResponse(rec), echo.New())
	token     interface{}
)

func TestCreate(t *testing.T) {
	_, err := tokenizer.Create(c, &a, &user)
	if err != nil {
		t.Error(err)
	}
	token = c.Get("jwt")
}

func TestParse(t *testing.T) {
	c.Request().Header.Set("Authorization", "Bearer "+token.(string))
	jwtParsed, err := tokenizer.Parse(c, &a)
	if err != nil {
		t.Error(err)
		return
	}
	if jwtParsed.Valid == false {
		t.Error("JWT not valid")
		return
	}
	if jwtParsed.Claims["email"] != user.Email || jwtParsed.Claims["name"] != user.Name {
		t.Error("Claims data are incorrect: " + jwtParsed.Claims["id"].(string) + "; " + jwtParsed.Claims["email"].(string) + "; " + jwtParsed.Claims["name"].(string))
	}
}
