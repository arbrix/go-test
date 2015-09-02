package user

// RegistrationForm is used when creating a user.
type RegistrationForm struct {
	Username string `json:"name" form:"registrationUsername" binding:"required"`
	Email    string `json:"email" form:"registrationEmail" binding:"required"`
	Password string `json:"pass" form:"registrationPassword" binding:"required"`
}

// LoginForm is used when creating a user authentication.
type LoginForm struct {
	Email    string `json:"login" form:"loginEmail" binding:"required"`
	Password string `json:"pass" form:"loginPassword" binding:"required"`
}

// UserForm is used when updating a user.
type UserForm struct {
	Age  int64  `json:"age" form:"age"`
	Name string `json:"name" form:"name" binding:"required"`
}

// PasswordForm is used when updating a user password.
type PasswordForm struct {
	CurrentPassword string `form:"currentPassword" binding:"required"`
	Password        string `form:"newPassword" binding:"required"`
}

// PasswordResetForm is used when reseting a password.
type PasswordResetForm struct {
	PasswordResetToken string `form:"token" binding:"required"`
	Password           string `form:"newPassword" binding:"required"`
}

// VerifyEmailForm is used when verifying an email.
type VerifyEmailForm struct {
	ActivationToken string `form:"token" binding:"required"`
}

// ActivateForm is used when activating user.
type ActivateForm struct {
	Activation bool `form:"activation" binding:"required"`
}

// UserRoleForm is used when adding or removing a role from a user.
type UserRoleForm struct {
	UserId int `form:"userId" binding:"required"`
	RoleId int `form:"roleId" binding:"required"`
}
