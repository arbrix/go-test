package oauth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/arbrix/go-test/interfaces"
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

type Facebook struct {
	oauth2.Config
	url        string
	RequestUrl *url.URL
	a          interfaces.App
}

func NewFacebook(a interfaces.App) (*Facebook, error) {
	var cIDVal, cScrVal, cRUrlVal, lAddrVal string
	var ok bool
	cID, err := a.GetConfig().Get("oauth-facebook-client-id")
	if err != nil {
		return nil, err
	}
	if cIDVal, ok = cID.(string); ok == false {
		return nil, errors.New("Client ID is not string type")
	}
	cScr, err := a.GetConfig().Get("oauth-facebook-client-secret")
	if err != nil {
		return nil, err
	}
	if cScrVal, ok = cScr.(string); ok == false {
		return nil, errors.New("Client Secret is not string type")
	}
	cRUrl, err := a.GetConfig().Get("oauth-facebook-redirect-url")
	if err != nil {
		return nil, err
	}
	if cRUrlVal, ok = cRUrl.(string); ok == false {
		return nil, errors.New("Client URL is not string type")
	}
	lAddr, err := a.GetConfig().Get("ListenAddress")
	if err != nil {
		return nil, err
	}
	if lAddrVal, ok = lAddr.(string); ok == false {
		return nil, errors.New("Server address is not string type")
	}

	fb := &Facebook{}
	fb.a = a

	fb.ClientID = cIDVal
	fb.ClientSecret = cScrVal
	fb.RedirectURL = "http://" + lAddrVal + cRUrlVal
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
	return fb, nil
}

// Your credentials should be obtained from the OAuth provider
func (fb *Facebook) URL() string {
	if len(fb.url) == 0 {
		fb.url = fb.AuthCodeURL("state")
	}
	return fb.url
}

func (fb *Facebook) Request(code string) (*http.Response, error) {
	var res *http.Response
	var req *http.Request
	var err error

	// Handle the exchange code to initiate a transport.
	token, err := fb.Exchange(oauth2.NoContext, code)
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
	tokenizer := jwt.NewTokenizer(fb.a)
	status, err := tokenizer.Create(c, user)
	if err != nil {
		return status, err
	}
	return http.StatusOK, nil
}

// CreateUser creates oauth user.
func (fb *Facebook) CreateUser(c *echo.Context, oauthUser *OauthUser) (*model.User, int, error) {
	var u *model.User
	tokenizer := jwt.NewTokenizer(fb.a)
	token, err := tokenizer.Parse(c)
	//TODO: add comma_ok statments
	if err == nil && token.Valid {
		u = &model.User{
			ID:    token.Claims["id"].(int64),
			Email: token.Claims["email"].(string),
			Name:  token.Claims["name"].(string),
		}
		return u, http.StatusBadRequest, errors.New("User already authorized")
	}
	u = &model.User{Email: oauthUser.Id + "@facebook.com", Name: oauthUser.Name}
	usrSrv := user.NewUserService(fb.a)
	usr, err := usrSrv.AddNew(u)
	if err != nil {
		return usr, http.StatusInternalServerError, errors.New("User is not created.")
	}
	return usr, http.StatusOK, nil
}

// LoginOrCreate login or create with oauthUser
func (fb *Facebook) LoginOrCreate(c *echo.Context, oauthUser *OauthUser) (int, error) {
	user := &model.User{}
	user.Email = oauthUser.Id + "@facebook.com"
	user.Name = oauthUser.Name
	if fb.a.GetDB().First(user, user) == nil {
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
	json.Unmarshal(body, &facebookUser)
	return *facebookUser, nil
}

// Oauth link connection and user.
func (fb *Facebook) Oauth(c *echo.Context) (int, error) {
	var code string
	var oauthUser OauthUser
	code = c.Form("code")
	response, err := fb.Request(code)
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
	return status, nil
}
