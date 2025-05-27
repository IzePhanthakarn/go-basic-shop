package usersHandlers

import (
	"strings"

	"github.com/IzePhanthakarn/kawaii-shop/config"
	"github.com/IzePhanthakarn/kawaii-shop/modules/entities"
	"github.com/IzePhanthakarn/kawaii-shop/modules/users"
	"github.com/IzePhanthakarn/kawaii-shop/modules/users/usersUsecases"
	"github.com/IzePhanthakarn/kawaii-shop/pkg/kawaiiauth"
	"github.com/gofiber/fiber/v3"
)

type userHandlersErrCode string

const (
	signUpCustomerErr     userHandlersErrCode = "users-001"
	signInErr             userHandlersErrCode = "users-002"
	refreshPassportErr    userHandlersErrCode = "users-003"
	signOutErr            userHandlersErrCode = "users-004"
	signUpAdminErr        userHandlersErrCode = "users-005"
	generateAdminTokenErr userHandlersErrCode = "users-006"
	getUserProfileErr     userHandlersErrCode = "users-007"
)

type IUserHandler interface {
	SignUpCustomer(c fiber.Ctx) error
	SignIn(c fiber.Ctx) error
	RefreshPassport(c fiber.Ctx) error
	SignOut(c fiber.Ctx) error
	SignUpAdmin(c fiber.Ctx) error
	GenerateAdminToken(c fiber.Ctx) error
	GetUserProfile(c fiber.Ctx) error
}

type userHandler struct {
	cfg          config.IConfig
	usersUsecase usersUsecases.IUsersUsecase
}

func UsersHandler(cfg config.IConfig, usersUsecase usersUsecases.IUsersUsecase) IUserHandler {
	return &userHandler{
		cfg:          cfg,
		usersUsecase: usersUsecase,
	}
}

func (h *userHandler) SignUpCustomer(c fiber.Ctx) error {
	// Request body parser
	req := new(users.UserRegisterReq)
	if err := c.Bind().Body(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signUpCustomerErr),
			err.Error(),
		).Res()
	}

	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signUpCustomerErr),
			"invalid email format",
		).Res()
	}

	// Insert user into database
	result, err := h.usersUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username already exists":
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		case "email already exists":
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.StatusInternalServerError,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		}
	}

	// Success response
	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *userHandler) SignUpAdmin(c fiber.Ctx) error {
	// Request body parser
	req := new(users.UserRegisterReq)
	if err := c.Bind().Body(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signUpCustomerErr),
			err.Error(),
		).Res()
	}

	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signUpCustomerErr),
			"invalid email format",
		).Res()
	}

	// Insert user into database
	result, err := h.usersUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username already exists":
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		case "email already exists":
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.StatusInternalServerError,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		}
	}

	// Success response
	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *userHandler) GenerateAdminToken(c fiber.Ctx) error {
	adminToken, err := kawaiiauth.NewKawaiiAuth(
		kawaiiauth.Admin,
		h.cfg.Jwt(),
		nil,
	)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(generateAdminTokenErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Token string `json:"token"`
		}{
			Token: adminToken.SignToken(),
		}).Res()
}

func (h *userHandler) SignIn(c fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.Bind().Body(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signInErr),
			err.Error(),
		).Res()
	}
	passport, err := h.usersUsecase.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signInErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *userHandler) RefreshPassport(c fiber.Ctx) error {
	req := new(users.UserRefreshCredential)
	if err := c.Bind().Body(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(refreshPassportErr),
			err.Error(),
		).Res()
	}
	passport, err := h.usersUsecase.RefreshPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(refreshPassportErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *userHandler) SignOut(c fiber.Ctx) error {
	req := new(users.UserRemoveCredential)
	if err := c.Bind().Body(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signOutErr),
			err.Error(),
		).Res()
	}
	err := h.usersUsecase.DeleteOauth(req.OauthId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(signOutErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}

func (h *userHandler) GetUserProfile(c fiber.Ctx) error {
	userId := strings.Trim(c.Params("user_id"), " ")

	// Get Profile
	result, err := h.usersUsecase.GetUserProfile(userId)
	if err != nil {
		switch err.Error() {
		case "get user failed: sql: no rows in result set":
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(getUserProfileErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.StatusInternalServerError,
				string(getUserProfileErr),
				err.Error(),
			).Res()
		}
	}

	// Success response
	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}
