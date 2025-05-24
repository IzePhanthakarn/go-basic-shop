package usersUsecases

import (
	"fmt"

	"github.com/IzePhanthakarn/kawaii-shop/config"
	"github.com/IzePhanthakarn/kawaii-shop/modules/users"
	"github.com/IzePhanthakarn/kawaii-shop/modules/users/usersRepositories"
	"github.com/IzePhanthakarn/kawaii-shop/pkg/kawaiiauth"
	"golang.org/x/crypto/bcrypt"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
	RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error)
	DeleteOauth(oauthId string) error
}

type usersUsecase struct {
	cfg             config.IConfig
	usersRepository usersRepositories.IUsersRepository
}

func UsersUsecase(cfg config.IConfig, usersRepository usersRepositories.IUsersRepository) IUsersUsecase {
	return &usersUsecase{
		cfg:             cfg,
		usersRepository: usersRepository,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	// Hash password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	// Insert user
	result, err := u.usersRepository.InsertUser(req, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *usersUsecase) InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error) {
	// Hash password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	// Insert user
	result, err := u.usersRepository.InsertUser(req, true)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *usersUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	// Find user
	user, err := u.usersRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// Sign token
	accessToken, err := kawaiiauth.NewKawaiiAuth(kawaiiauth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id:       user.Id,
		RoleId:   user.RoleId,
	})
	if err != nil {
		return nil, err
	}
	refreshToken, err := kawaiiauth.NewKawaiiAuth(kawaiiauth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id:       user.Id,
		RoleId:   user.RoleId,
	})
	if err != nil {
		return nil, err
	}

	// Set passport
	passport := &users.UserPassport{
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken: accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}
	
	// Insert oauth
	if err := u.usersRepository.InsertOauth(passport); err != nil {
		return nil, err
	}

	return passport, nil
}

func (u *usersUsecase) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
	// Parse token
	claims, err := kawaiiauth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Find oauth
	oauth, err := u.usersRepository.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Find user
	user, err := u.usersRepository.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	newClaims := &users.UserClaims{
		Id:       user.Id,
		RoleId:   user.RoleId,
	}

	// Sign token
	accessToken, err := kawaiiauth.NewKawaiiAuth(kawaiiauth.Access, u.cfg.Jwt(), newClaims)
	if err != nil {
		return nil, err
	}
	
	refreshToken := kawaiiauth.RepeatToken(
		u.cfg.Jwt(),
		newClaims,
		claims.ExpiresAt.Unix(),
	)

	passport := &users.UserPassport{
		User: user,
		Token: &users.UserToken{
			Id: oauth.Id,
			AccessToken: accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}
	
	// Update oauth
	if err := u.usersRepository.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}
	
	return passport, nil
}

func (u *usersUsecase) DeleteOauth(oauthId string) error {
	return u.usersRepository.DeleteOauth(oauthId)
}