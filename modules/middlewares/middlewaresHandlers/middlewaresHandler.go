package middlewaresHandlers

import (
	"strings"

	"github.com/IzePhanthakarn/go-basic-shop/config"
	"github.com/IzePhanthakarn/go-basic-shop/modules/entities"
	"github.com/IzePhanthakarn/go-basic-shop/modules/middlewares/middlewaresUsecases"
	"github.com/IzePhanthakarn/go-basic-shop/pkg/auth"
	"github.com/IzePhanthakarn/go-basic-shop/pkg/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

type middlewareHandlersErrCode string

const (
	routerCheckErr middlewareHandlersErrCode = "middleware-001"
	jwtAuthErr     middlewareHandlersErrCode = "middleware-002"
	paramsCheckErr middlewareHandlersErrCode = "middleware-003"
	authorizeError middlewareHandlersErrCode = "middleware-004"
	apiKeyAuthErr  middlewareHandlersErrCode = "middleware-005"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
	Authorize(expectRoleId ...int) fiber.Handler
	ApiKeyAuth() fiber.Handler
}

type middlewaresHandler struct {
	cfg               config.IConfig
	middlewareUsecase middlewaresUsecases.IMiddlewaresUsecase
}

func MiddlewaresHandler(cfg config.IConfig, middlewareUsecase middlewaresUsecases.IMiddlewaresUsecase) IMiddlewaresHandler {
	return &middlewaresHandler{
		cfg:               cfg,
		middlewareUsecase: middlewareUsecase,
	}
}

func (m *middlewaresHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"},
		AllowHeaders:     []string{""},
		AllowCredentials: false,
		ExposeHeaders:    []string{""},
		MaxAge:           0,
	})
}

func (h *middlewaresHandler) RouterCheck() fiber.Handler {
	return func(c fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"router not found",
		).Res()
	}
}

func (h *middlewaresHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006 15:04:05",
		TimeZone:   "Asia/Bangkok",
	})
}

func (h *middlewaresHandler) JwtAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := auth.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.StatusUnauthorized,
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}

		claims := result.Claims
		if !h.middlewareUsecase.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(
				fiber.StatusUnauthorized,
				string(jwtAuthErr),
				"invalid access token",
			).Res()
		}

		// Set UserId
		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)

		return c.Next()
	}
}

func (h *middlewaresHandler) ParamsCheck() fiber.Handler {
	return func(c fiber.Ctx) error {
		userId := c.Locals("userId")
		if c.Locals("userRoleId").(int) == 2 {
			return c.Next()
		}
		if c.Params("user_id") != userId {
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(paramsCheckErr),
				"invalid user id",
			).Res()
		}
		return c.Next()
	}
}

func (h *middlewaresHandler) Authorize(expectRoleId ...int) fiber.Handler {
	return func(c fiber.Ctx) error {
		userRoleId, ok := c.Locals("userRoleId").(int)
		if !ok {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(authorizeError),
				"invalid user role id",
			).Res()
		}

		roles, err := h.middlewareUsecase.FindRole()
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(authorizeError),
				err.Error(),
			).Res()
		}

		sum := 0
		for _, v := range expectRoleId {
			sum += v
		}

		expectedValueBinary := utils.BinaryConverter(sum, len(roles))
		userValueBinary := utils.BinaryConverter(userRoleId, len(roles))

		for i := range userValueBinary {
			if userValueBinary[i]&expectedValueBinary[i] == 1 {
				return c.Next()
			}
		}

		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(authorizeError),
			"unauthorized",
		).Res()
	}
}

func (h *middlewaresHandler) ApiKeyAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		key := c.Get("x-api-key")
		if _, err := auth.ParseApiKey(h.cfg.Jwt(), key); err != nil {
			return entities.NewResponse(c).Error(
				fiber.StatusUnauthorized,
				string(apiKeyAuthErr),
				"invalid api key",
			).Res()
		}
		return c.Next()
	}
}
