package appinfoHandlers

import (
	"github.com/IzePhanthakarn/kawaii-shop/config"
	"github.com/IzePhanthakarn/kawaii-shop/modules/appinfo/appinfoUsecases"
	"github.com/IzePhanthakarn/kawaii-shop/modules/entities"
	"github.com/IzePhanthakarn/kawaii-shop/pkg/kawaiiauth"
	"github.com/gofiber/fiber/v3"
)


type appinfoHandlersErrCode string

const (
	generateApiKeyErr appinfoHandlersErrCode = "appinfo-001"
)

type IAppinfoHandler interface {
	GenerateApiKey(c fiber.Ctx) error
}

type appinfoHandler struct {
	cfg             config.IConfig
	appinfoUsecases appinfoUsecases.IAppinfoUsecase
}

func AppinfoHandler(cfg config.IConfig, appinfoUsecases appinfoUsecases.IAppinfoUsecase) IAppinfoHandler {
	return &appinfoHandler{
		cfg:             cfg,
		appinfoUsecases: appinfoUsecases,
	}
}

func (h *appinfoHandler) GenerateApiKey(c fiber.Ctx) error {
	apiKey, err := kawaiiauth.NewKawaiiAuth(
		kawaiiauth.ApiKey,
		h.cfg.Jwt(),
		nil,
	)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(generateApiKeyErr),
			err.Error(),
		).Res()
	}
		
	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			ApiKey string `json:"api_key"`
		}{
			ApiKey: apiKey.SignToken(),
		}).Res()
}
