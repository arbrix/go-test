package model

import (
	"time"
)

type omit bool

// User is a user model
type User struct {
	Id              int64     `json:"id"`
	Email           string    `json:"email",sql:"size:255;unique"`
	Password        string    `json:"password",sql:"size:255"`
	Name            string    `json:"name",sql:"size:255"`
	Username        string    `json:"username",sql:"size:127;unique"`
	Birthday        time.Time `json:"birthday"`
	Gender          int8      `json:"gender"`
	Description     string    `json:"description",sql:"size:127"`
	Token           string    `json:"token"`
	TokenExpiration time.Time `json:"tokenExperiation"`
	// admin
	Activation         bool      `json:"activation"`
	PasswordResetToken string    `json:"passwordResetToken"`
	ActivationToken    string    `json:"activationToken"`
	PasswordResetUntil time.Time `json:"passwordResetUntil"`
	ActivateUntil      time.Time `json:"activateUntil"`
	ActivatedAt        time.Time `json:"activatedAt"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	DeletedAt          time.Time `json:"deletedAt"`
	LastLoginAt        time.Time `json:"lastLoginAt"`
	CurrentLoginAt     time.Time `json:"currentLoginAt"`
	LastLoginIp        string    `json:"lastLoginIp",sql:"size:100"`
	CurrentLoginIp     string    `json:"currentLoginIp",sql:"size:100"`
}

// PublicUser is a public user model that contains only a few information for everyone.
type PublicUser struct {
	*User
	Email           omit `json:"email,omitempty",sql:"size:255;unique"`
	Password        omit `json:"password,omitempty",sql:"size:255"`
	Name            omit `json:"name,omitempty",sql:"size:255"`
	Birthday        omit `json:"birthday,omitempty"`
	Gender          omit `json:"gender,omitempty"`
	Token           omit `json:"token,omitempty"`
	TokenExpiration omit `json:"tokenExperiation,omitempty"`

	// admin
	Activation         omit `json:"activation,omitempty"`
	PasswordResetToken omit `json:"passwordResetToken,omitempty"`
	ActivationToken    omit `json:"activationToken,omitempty"`
	PasswordResetUntil omit `json:"passwordResetUntil,omitempty"`
	ActivateUntil      omit `json:"activateUntil,omitempty"`
	ActivatedAt        omit `json:"activatedAt,omitempty"`
	UpdatedAt          omit `json:"updatedAt,omitempty"`
	DeletedAt          omit `json:"deletedAt,omitempty"`
	LastLoginAt        omit `json:"lastLoginAt,omitempty"`
	CurrentLoginAt     omit `json:"currentLoginAt,omitempty"`
	LastLoginIp        omit `json:"lastLoginIp,omitempty",sql:"size:100"`
	CurrentLoginIp     omit `json:"currentLoginIp,omitempty",sql:"size:100"`
}

// Connection is a connection model for oauth.
type Connection struct {
	Id             int64  `json:"id"`
	UserId         int64  `json:"userId"`
	ProviderId     int64  `gorm:"column:provider_id", json:"providerId"`
	ProviderUserId string `gorm:"column:provider_user_id", json:"providerUserId"`
	AccessToken    string `json:"accessToken"`
	ProfileUrl     string `gorm:"column:profile_url", json:"profileUrl"`
	ImageUrl       string `gorm:"column:image_url", json:"imageUrl"`
}
