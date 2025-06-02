package appinfoHandlers

import (
	"strconv"
	"strings"

	"github.com/IzePhanthakarn/kawaii-shop/config"
	"github.com/IzePhanthakarn/kawaii-shop/modules/appinfo"
	"github.com/IzePhanthakarn/kawaii-shop/modules/appinfo/appinfoUsecases"
	"github.com/IzePhanthakarn/kawaii-shop/modules/entities"
	"github.com/IzePhanthakarn/kawaii-shop/pkg/kawaiiauth"
	"github.com/gofiber/fiber/v3"
)

type appinfoHandlersErrCode string

const (
	generateApiKeyErr appinfoHandlersErrCode = "appinfo-001"
	findCategoryErr   appinfoHandlersErrCode = "appinfo-002"
	addCategoryErr    appinfoHandlersErrCode = "appinfo-003"
	removeCategoryErr appinfoHandlersErrCode = "appinfo-004"
)

type IAppinfoHandler interface {
	GenerateApiKey(c fiber.Ctx) error
	FindCategory(c fiber.Ctx) error
	AddCategory(c fiber.Ctx) error
	RemoveCategory(c fiber.Ctx) error
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

// @Summary Generate Api Key
// @Description Generate Api Key
// @Tags Appinfo
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} appinfo.GenerateApiKeyRes
// @Router /appinfo/apikey [get]
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

// @Summary Find Categories
// @Description Find Categories
// @Tags Categories
// @Accept  json
// @Produce  json
// @Param title path string true "Title"
// @Security BearerAuth
// @Success 200 {array} appinfo.Category
// @Router /appinfo/categories/{title} [get]
func (h *appinfoHandler) FindCategory(c fiber.Ctx) error {
	req := new(appinfo.CategoryFilter)
	if err := c.Bind().Query(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}

	category, err := h.appinfoUsecases.FindCategory(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		category,
	).Res()
}

// @Summary Add Caregory
// @Description Add Caregory
// @Tags Categories
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body appinfo.Category true "Category Request"
// @Success 200 {array} appinfo.Category
// @Router /appinfo/categories [post]
func (h *appinfoHandler) AddCategory(c fiber.Ctx) error {
	req := make([]*appinfo.Category, 0)
	if err := c.Bind().Body(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(addCategoryErr),
			err.Error(),
		).Res()
	}

	if len(req) == 0 {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(addCategoryErr),
			"request body is empty",
		).Res()
	}

	if err := h.appinfoUsecases.InsertCategory(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(addCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusCreated,
		req,
	).Res()
}

// @Summary Delete File
// @Description Delete File
// @Tags Files
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param category_id path string true "Category Id"
// @Success 200 {object} appinfo.CategoryRemoveRes
// @Router /appinfo/categories/{category_id} [delete]
func (h *appinfoHandler) RemoveCategory(c fiber.Ctx) error {
	categoryId := strings.Trim(c.Params("category_id"), " ")
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(removeCategoryErr),
			"invalid category_id",
		).Res()
	}

	if categoryIdInt <= 0 {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(removeCategoryErr),
			"category_id must be greater than 0",
		).Res()
	}

	if err := h.appinfoUsecases.DeleteCategory(categoryIdInt); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(removeCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			CategoryId int `json:"category_id"`
		}{
			CategoryId: categoryIdInt,
		},
	).Res()
}
