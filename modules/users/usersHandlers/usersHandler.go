package usersHandlers

import (
	"github.com/IzePhanthakarn/kawaii-shop/config"
	"github.com/IzePhanthakarn/kawaii-shop/modules/entities"
	"github.com/IzePhanthakarn/kawaii-shop/modules/users"
	"github.com/IzePhanthakarn/kawaii-shop/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v3"
)

type userHandlersErrCode string

const (
	signUpCustomerErr userHandlersErrCode = "users-001"
	signInErr         userHandlersErrCode = "users-002"
)

type IUserHandler interface {
	SignUpCustomer(c fiber.Ctx) error
	SignIn(c fiber.Ctx) error
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