package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/arbrix/go-test/app"
	"github.com/arbrix/go-test/model"
	"github.com/arbrix/go-test/service/user"
	"github.com/arbrix/go-test/util/jwt"
	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

// OauthUser is a struct for connection.
type OauthUser struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type AuthResponse struct {
	State        string `form:"state"`
	Code         string `form:"code"`
	Authuser     int    `form:"authuser"`
	NumSessions  int    `form:"num_sessions"`
	prompt       string `form:"prompt"`
	ImageName    string `form:"imageName"`
	SessionState string `form:"session_state"`
}

type Facebook struct {
	oauth2.Config
	url        string
	RequestUrl *url.URL
	a          *app.App
}

func (fb *Facebook) Init(a *app.App) error {
	cID, err := a.GetConfig().Get("oauth-facebook-client-id")
	if err != nil {
		return err
	}
	cScr, err := a.GetConfig().Get("oauth-facebook-client-secret")
	if err != nil {
		return err
	}
	cRUrl, err := a.GetConfig().Get("oauth-facebook-redirect-url")
	if err != nil {
		return err
	}
	lAddr, err := a.GetConfig().Get("ListenAddress")
	if err != nil {
		return err
	}

	fb.a = a
	fb.ClientID = cID.(string)
	fb.ClientSecret = cScr.(string)
	fb.RedirectURL = "http://" + lAddr.(string) + cRUrl.(string)
	fb.Scopes = []string{
		"public_profile",
		"email",
	}
	fb.Endpoint = facebook.Endpoint
	fb.RequestUrl = &url.URL{
		Scheme: "https",
		Host:   "graph.facebook.com",
		Opaque: "//graph.facebook.com/me",
	}
	return nil
}

// Your credentials should be obtained from the OAuth provider
func (fb *Facebook) URL() string {
	if len(fb.url) == 0 {
		fb.url = fb.AuthCodeURL("state")
	}
	return fb.url
}

func (fb *Facebook) Request(authResponse AuthResponse) (*http.Response, error) {
	var res *http.Response
	var req *http.Request
	var err error

	// Handle the exchange code to initiate a transport.
	token, err := fb.Exchange(oauth2.NoContext, authResponse.Code)
	if err != nil {
		return res, err
	}
	client := fb.Client(oauth2.NoContext, token)
	req, err = http.NewRequest("GET", "url.String()", nil)
	if err != nil {
		return res, err
	}
	req.URL = fb.RequestUrl
	res, err = client.Do(req)
	return res, nil
}

// Login login with oauthUser's username.
func (fb *Facebook) Login(c *echo.Context, user *model.User) (int, error) {
	tokenizer := jwt.Token{}
	status, err := tokenizer.Create(c, fb.a, user)
	if err != nil {
		return status, err
	}
	return http.StatusOK, nil
}

// CreateUser creates oauth user.
func (fb *Facebook) CreateUser(c *echo.Context, oauthUser *OauthUser) (*model.User, int, error) {
	var u model.User
	tokenizer := jwt.Token{}
	token, err := tokenizer.Parse(c, fb.a)
	if err == nil && token.Valid {
		u = model.User{
			ID:    token.Claims["id"].(int64),
			Email: token.Claims["email"].(string),
			Name:  token.Claims["name"].(string),
		}
		return &u, http.StatusBadRequest, errors.New("User already authorized")
	}
	u = model.User{Email: oauthUser.Email, Name: oauthUser.Name}
	if len(u.Name) == 0 {
		if len(u.Email) > 0 {
			u.Name = strings.Split(u.Email, "@")[0]
		} else {
			u.Name = "OauthUser"
		}
	}
	usrSrv := user.Service{}
	usr, err := usrSrv.AddNew(&u, fb.a)
	if err != nil {
		return usr, http.StatusInternalServerError, errors.New("User is not created.")
	}
	return usr, http.StatusOK, nil
}

// LoginOrCreate login or create with oauthUser
func (fb *Facebook) LoginOrCreate(c *echo.Context, oauthUser *OauthUser) (int, error) {
	var user *model.User

	if fb.a.GetDB().First(user, map[string]interface{}{"email": oauthUser.Email, "name": oauthUser.Name}) == nil {
		return fb.Login(c, user)
	}

	user, status, err := fb.CreateUser(c, oauthUser)
	if err != nil {
		return status, err
	}
	return fb.Login(c, user)
}

// SetUser set facebook user.
func (fb *Facebook) SetUser(response *http.Response) (OauthUser, error) {
	facebookUser := &OauthUser{}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return *facebookUser, err
	}
	fmt.Println("%s", body) //TODO: delete after tests
	json.Unmarshal(body, &facebookUser)
	return *facebookUser, nil
}

// Oauth link connection and user.
func (fb *Facebook) Oauth(c *echo.Context) (int, error) {
	var authResponse AuthResponse
	var oauthUser OauthUser
	err := c.Bind(&authResponse)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	response, err := fb.Request(authResponse)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	oauthUser, err = fb.SetUser(response)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	status, err := fb.LoginOrCreate(c, &oauthUser)
	if err != nil {
		return status, err
	}
	return http.StatusSeeOther, nil
}
