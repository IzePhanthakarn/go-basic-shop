package users

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       string `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Username string `db:"username" json:"username"`
	RoleId   int    `db:"role_id" json:"role_id"`
}

type UserRegisterReq struct {
	Email    string `db:"email" json:"email" form:"email"`
	Username string `db:"username" json:"username" form:"username"`
	Password string `db:"password" json:"password" form:"password"`
}

type UserCredential struct {
	Email    string `db:"email" json:"email" form:"email"`
	Password string `db:"password" json:"password" form:"password"`
}

type UserCredentialCheck struct {
	Id       string `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Username string `db:"username"`
	RoleId   int    `db:"role_id"`
}

type AdminTokenResponse struct {
    Token string `json:"token"`
}

func (obj *UserRegisterReq) BcryptHashing() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	obj.Password = string(hashedPassword)
	return nil
}

func (obj *UserRegisterReq) IsEmail() bool {
	match, err := regexp.Match(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, []byte(obj.Email))
	if err != nil {
		return false
	}

	return match
}

type UserPassport struct {
	User  *User      `json:"user"`
	Token *UserToken `json:"token"`
}

type UserToken struct {
	Id           string `id:"id" json:"id"`
	AccessToken  string `db:"access_token" json:"access_token"`
	RefreshToken string `db:"refresh_token" json:"refresh_token"`
}

type UserClaims struct {
	Id string `db:"id" json:"id"`
	RoleId int `db:"role" json:"role"`
}

type UserRefreshCredential struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
}

type Oauth struct {
	Id string `db:"id" json:"id"`
	UserId string `db:"user_id" json:"user_id"`
}

type UserRemoveCredential struct {
	OauthId string `json:"oauth_id" form:"oauth_id"`
}