package facebook

import (
	"fmt"
	"net/url"

	"github.com/arbrix/go-test/config"
	"golang.org/x/oauth2"
)

const (
	ProviderId = 4
	Scheme     = "https"
	Host       = "graph.facebook.com"
	Opaque     = "//graph.facebook.com/me"
	AuthURL    = "https://www.facebook.com/dialog/oauth"
	TokenURL   = "https://graph.facebook.com/oauth/access_token"
)

var RequestURL = &url.URL{
	Scheme: Scheme,
	Host:   Host,
	Opaque: Opaque,
}

var Endpoint = oauth2.Endpoint{
	AuthURL:  AuthURL,
	TokenURL: TokenURL,
}

var Config = Oauth2Config()

func Oauth2Config() *oauth2.Config {
	fmt.Println("http://" + config.JsonConfig.ListenAddress + config.OauthFacebookRedirectURL)
	return &oauth2.Config{
		ClientID:     config.OauthFacebookClientID,
		ClientSecret: config.OauthFacebookClientSecret,
		RedirectURL:  "http://" + config.JsonConfig.ListenAddress + config.OauthFacebookRedirectURL,
		Scopes: []string{
			"public_profile",
			"email",
		},
		Endpoint: Endpoint,
	}
}
