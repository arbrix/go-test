package user

import (
	"errors"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/arbrix/go-test/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var userFields = []string{"name", "email", "createdAt", "updatedAt"}

// SuggestUsername suggest user's name if user's name already occupied.
func SuggestUsername(username string) string {
	var count int
	var usernameCandidate string
	db.ORM.Model(model.User{}).Where(&model.User{Username: username}).Count(&count)
	log.Debugf("count Before : %d", count)
	if count == 0 {
		return username
	} else {
		var postfix int
		for {
			usernameCandidate = username + strconv.Itoa(postfix)
			log.Debugf("usernameCandidate: %s\n", usernameCandidate)
			db.ORM.Model(model.User{}).Where(&model.User{Username: usernameCandidate}).Count(&count)
			log.Debugf("count after : %d\n", count)
			postfix = postfix + 1
			if count == 0 {
				break
			}
		}
	}
	return usernameCandidate
}

// CreateUserFromJson creates a user from a registration form.
func CreateUserFromJson(registrationForm RegistrationForm) (model.User, error) {
	var user model.User
	log.Debugf("registrationForm %+v\n", registrationForm)
	modelHelper.AssignValue(&user, &registrationForm)
	token, err := crypto.GenerateRandomToken32()
	if err != nil {
		return user, errors.New("Token not generated.")
	}
	user.Token = token
	user.TokenExpiration = timeHelper.FewDaysLater(config.AuthTokenExpirationDay)
	log.Debugf("user %+v\n", user)
	if db.ORM.Create(&user).Error != nil {
		return user, errors.New("User is not created.")
	}
	return user, nil
}

// CreateUser creates a user.
func CreateUser(c *gin.Context) (int, error) {
	var registrationForm RegistrationForm
	var user model.User
	var status int
	var err error

	err = c.BindWith(&registrationForm, binding.JSON)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	userPass := registrationForm.Password
	password, err := bcrypt.GenerateFromPassword([]byte(registrationForm.Password), 10)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	registrationForm.Password = string(password)
	user, err = CreateUserFromJson(registrationForm)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	status, err = Login(user.Email, userPass)
	return status, err
}

// RetrieveUser retrieves a user.
func RetrieveUser(c *gin.Context) (*model.PublicUser, bool, int64, int, error) {
	var user model.User
	var currentUserId int64
	var isAuthor bool
	// var publicUser *model.PublicUser
	// publicUser.User = &user
	id := c.Params.ByName("id")
	if db.ORM.Select(config.UserPublicFields).First(&user, id).RecordNotFound() {
		return &model.PublicUser{User: &user}, isAuthor, currentUserId, http.StatusNotFound, errors.New("User is not found.")
	}
	currentUser, err := CurrentUser(c)
	if err == nil {
		currentUserId = currentUser.Id
		isAuthor = currentUser.Id == user.Id
	}
	return &model.PublicUser{User: &user}, isAuthor, currentUserId, http.StatusOK, nil
}

// RetrieveUsers retrieves users.
func RetrieveUsers(c *gin.Context) []*model.PublicUser {
	var users []*model.User
	var userArr []*model.PublicUser
	db.ORM.Select(config.UserPublicFields).Find(&users)
	for _, user := range users {
		userArr = append(userArr, &model.PublicUser{User: user})
	}
	return userArr
}

// UpdateUserCore updates a user. (Applying the modifed data of user).
func UpdateUserCore(user *model.User) (int, error) {
	token, err := crypto.GenerateRandomToken32()
	if err != nil {
		return http.StatusInternalServerError, errors.New("Token not generated.")
	}
	user.Token = token
	user.TokenExpiration = timeHelper.FewDaysLater(config.AuthTokenExpirationDay)
	if db.ORM.Save(user).Error != nil {
		return http.StatusInternalServerError, errors.New("User is not updated.")
	}
	return http.StatusOK, nil
}

// UpdateUser updates a user.
func UpdateUser(c *gin.Context) (*model.User, int, error) {
	id := c.Params.ByName("id")
	var user model.User
	if db.ORM.First(&user, id).RecordNotFound() {
		return &user, http.StatusNotFound, errors.New("User is not found.")
	}
	switch c.Request.FormValue("type") {
	case "password":
		var passwordForm PasswordForm
		c.BindWith(&passwordForm, binding.Form)
		log.Debugf("form %+v\n", passwordForm)
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordForm.CurrentPassword))
		if err != nil {
			log.Error("Password Incorrect.")
			return &user, http.StatusInternalServerError, errors.New("User is not updated. Password Incorrect.")
		} else {
			newPassword, err := bcrypt.GenerateFromPassword([]byte(passwordForm.Password), 10)
			if err != nil {
				return &user, http.StatusInternalServerError, errors.New("User is not updated. Password not Generated.")
			} else {
				passwordForm.Password = string(newPassword)
				modelHelper.AssignValue(&user, &passwordForm)
			}
		}
	default:
		var form UserForm
		c.BindWith(&form, binding.Form)
		log.Debugf("form %+v\n", form)
		modelHelper.AssignValue(&user, &form)
	}

	log.Debugf("params %+v\n", c.Params)
	status, err := UpdateUserCore(&user)
	if err != nil {
		return &user, status, err
	}
	return &user, status, err
}

// DeleteUser deletes a user.
func DeleteUser(c *gin.Context) (int, error) {
	id := c.Params.ByName("id")
	var user model.User
	if db.ORM.First(&user, id).RecordNotFound() {
		return http.StatusNotFound, errors.New("User is not found.")
	}
	if db.ORM.Delete(&user).Error != nil {
		return http.StatusInternalServerError, errors.New("User is not deleted.")
	}
	return http.StatusOK, nil
}

// RetrieveCurrentUser retrieves a current user.
func RetrieveCurrentUser(c *gin.Context) (model.User, int, error) {
	user, err := CurrentUser(c)
	if err != nil {
		return user, http.StatusInternalServerError, err
	}
	return user, http.StatusOK, nil
}

// RetrieveUserByEmail retrieves a user by an email
func RetrieveUserByEmail(c *gin.Context) (*model.PublicUser, string, int, error) {
	email := c.Params.ByName("email")
	var user model.User
	if db.ORM.Unscoped().Select(config.UserPublicFields).Where("email like ?", "%"+email+"%").First(&user).RecordNotFound() {
		return &model.PublicUser{User: &user}, email, http.StatusNotFound, errors.New("User is not found.")
	}
	return &model.PublicUser{User: &user}, email, http.StatusOK, nil
}

// RetrieveUsersByEmail retrieves users by an email
func RetrieveUsersByEmail(c *gin.Context) []*model.PublicUser {
	var users []*model.User
	var userArr []*model.PublicUser
	email := c.Params.ByName("email")
	db.ORM.Select(config.UserPublicFields).Where("email like ?", "%"+email+"%").Find(&users)
	for _, user := range users {
		userArr = append(userArr, &model.PublicUser{User: user})
	}
	return userArr
}

// RetrieveUserByUsername retrieves a user by username.
func RetrieveUserByUsername(c *gin.Context) (*model.PublicUser, string, int, error) {
	username := c.Params.ByName("username")
	var user model.User
	if db.ORM.Unscoped().Select(config.UserPublicFields).Where("username like ?", "%"+username+"%").First(&user).RecordNotFound() {
		return &model.PublicUser{User: &user}, username, http.StatusNotFound, errors.New("User is not found.")
	}
	return &model.PublicUser{User: &user}, username, http.StatusOK, nil
}

// RetrieveUserForAdmin retrieves a user for an administrator.
func RetrieveUserForAdmin(c *gin.Context) (model.User, int, error) {
	id := c.Params.ByName("id")
	var user model.User
	if db.ORM.First(&user, id).RecordNotFound() {
		return user, http.StatusNotFound, errors.New("User is not found.")
	}
	return user, http.StatusOK, nil
}

// RetrieveUsersForAdmin retrieves users for an administrator.
func RetrieveUsersForAdmin(c *gin.Context) []model.User {
	var users []model.User
	db.ORM.Find(&users)
	return users
}

// ActivateUser toggle activation of a user.
func ActivateUser(c *gin.Context) (model.User, int, error) {
	id := c.Params.ByName("id")
	var user model.User
	var form ActivateForm
	c.BindWith(&form, binding.Form)
	if db.ORM.First(&user, id).RecordNotFound() {
		return user, http.StatusNotFound, errors.New("User is not found.")
	}
	user.Activation = form.Activation
	if db.ORM.Save(&user).Error != nil {
		return user, http.StatusInternalServerError, errors.New("User not activated.")
	}
	return user, http.StatusOK, nil
}
