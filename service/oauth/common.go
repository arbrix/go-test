package oauth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/arbrix/go-test/app"
	"github.com/arbrix/go-test/model"
	"github.com/arbrix/go-test/service/user"
	"github.com/arbrix/go-test/util/oauth2/facebook"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

type oauthStatusMap map[string]bool

// OauthUser is a struct for connection.
type OauthUser struct {
	Id         string
	Email      string
	Username   string
	Name       string
	ImageUrl   string
	ProfileUrl string
}

// getLoginName get login name from user. It is user's username or email.
// func getLoginName(user model.User) string {
// 	username := user.Username
// 	loginName := user.Email
// 	if len(username) > 0 {
// 		loginName = username
// 	}
// 	return loginName
// }

// RetrieveOauthStatus retrieves ouath status that provider connected or not.
func RetrieveOauthStatus(c *gin.Context) (oauthStatusMap, int, error) {
	var oauthStatus oauthStatusMap
	var connections []model.Connection
	oauthStatus = make(oauthStatusMap)
	currentUser, err := user.CurrentUser(c)
	if err != nil {
		return oauthStatus, http.StatusUnauthorized, err
	}
	db.ORM.Model(&currentUser).Association("Connections").Find(&connections)
	for _, connection := range connections {
		log.Debugf("connection.ProviderId : %d", connection.ProviderId)
		switch connection.ProviderId {
		case facebook.ProviderId:
			oauthStatus["facebook"] = true
		}
	}
	log.Debugf("oauthStatus : %v", oauthStatus)
	return oauthStatus, http.StatusOK, nil
}

// LoginWithOauthUser login with oauthUser's username.
func LoginWithOauthUser(c *gin.Context, userContact string) (int, error) {
	status, err := user.CreateToken(userContact)
	if err != nil {
		return status, err
	}
	return http.StatusOK, nil
}

// CreateOauthUser creates oauth user.
func CreateOauthUser(c *gin.Context, oauthUser *OauthUser, connection *model.Connection) (model.User, int, error) {
	var registrationForm user.RegistrationForm
	var user model.User
	modelHelper.AssignValue(&registrationForm, oauthUser)
	registrationForm.Password = random.GenerateRandomString(12)
	if len(registrationForm.Username) == 0 {
		if len(registrationForm.Email) > 0 {
			registrationForm.Username = strings.Split(registrationForm.Email, "@")[0]
		} else {
			registrationForm.Username = "OauthUser"
		}
	}
	registrationForm.Username = user.SuggestUsername(registrationForm.Username)
	generatedPassword, err := bcrypt.GenerateFromPassword([]byte(registrationForm.Password), 10)
	if err != nil {
		return user, http.StatusInternalServerError, errors.New("Password not generated.")
	}
	registrationForm.Password = string(generatedPassword)
	currentUser, err := user.CurrentUser(c)
	if err != nil {
		log.Errorf("currentUser : %v", currentUser)
		user, err = user.CreateUserFromJson(registrationForm)
		if err != nil {
			return user, http.StatusInternalServerError, errors.New("User is not created.")
		}
	} else {
		if db.ORM.Where("id = ?", currentUser.Id).First(&user).RecordNotFound() {
			return user, http.StatusInternalServerError, errors.New("User is not found.")
		}
	}
	return user, http.StatusOK, nil
}

// LoginOrCreateOauthUser login or create with oauthUser
func LoginOrCreateOauthUser(c *gin.Context, oauthUser *OauthUser, providerID int64, token *oauth2.Token) (int, error) {
	var connection model.Connection
	var count int
	db.ORM.Where("provider_id = ? and provider_user_id = ?", providerID, oauthUser.Id).First(&connection).Count(&count)

	connection.ProviderId = providerID
	connection.ProviderUserId = oauthUser.Id
	connection.AccessToken = token.AccessToken
	connection.ProfileUrl = oauthUser.ProfileUrl
	connection.ImageUrl = oauthUser.ImageUrl
	log.Debugf("connection count : %v", count)
	var user model.User
	if count == 1 {
		if db.ORM.First(&user, "id = ?", connection.UserId).RecordNotFound() {
			return http.StatusNotFound, errors.New("User is not found.")
		}
		log.Debugf("user : %v", user)
		if db.ORM.Save(&connection).Error != nil {
			return http.StatusInternalServerError, errors.New("Connection is not updated.")
		}
		status, err := LoginWithOauthUser(c, user.Email)
		return status, err
	}
	log.Debugf("Connection is not exist.")
	_, status, err := CreateOauthUser(c, oauthUser, &connection)
	if err != nil {
		return status, err
	}
	status, err = LoginWithOauthUser(c, user.Email)
	return status, err
}

// RevokeOauth revokes oauth connection.
func RevokeOauth(c *gin.Context, providerID int64) (oauthStatusMap, int, error) {
	var oauthStatus oauthStatusMap
	var connection model.Connection
	currentUser, err := user.CurrentUser(c)
	if err != nil {
		return oauthStatus, http.StatusUnauthorized, err
	}
	if db.ORM.First(&connection, "user_id= ? and provider_id = ? ", currentUser.Id, providerID).RecordNotFound() {
		return oauthStatus, http.StatusNotFound, errors.New("Connection is not found.")
	}
	if db.ORM.Delete(&connection).Error != nil {
		return oauthStatus, http.StatusInternalServerError, errors.New("Connection not revoked from user.")
	}
	oauthStatus, status, err := RetrieveOauthStatus(c)
	return oauthStatus, status, err
}
