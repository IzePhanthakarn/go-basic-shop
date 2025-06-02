package productsHandlers

import (
	"fmt"
	"strings"

	"github.com/IzePhanthakarn/kawaii-shop/config"
	"github.com/IzePhanthakarn/kawaii-shop/modules/appinfo"
	"github.com/IzePhanthakarn/kawaii-shop/modules/entities"
	"github.com/IzePhanthakarn/kawaii-shop/modules/files"
	"github.com/IzePhanthakarn/kawaii-shop/modules/files/filesUsecases"
	"github.com/IzePhanthakarn/kawaii-shop/modules/products"
	"github.com/IzePhanthakarn/kawaii-shop/modules/products/productsUsecases"
	"github.com/gofiber/fiber/v3"
)

type productsHanflersErrCode string

const (
	findOneProductErr productsHanflersErrCode = "products-001"
	findProductErr    productsHanflersErrCode = "products-002"
	insertProductErr  productsHanflersErrCode = "products-003"
	updateProductErr  productsHanflersErrCode = "products-004"
	deleteProductErr  productsHanflersErrCode = "products-005"
)

type IProductsHandler interface {
	FindOneProduct(c fiber.Ctx) error
	FindProduct(c fiber.Ctx) error
	AddProduct(c fiber.Ctx) error
	UpdateProduct(c fiber.Ctx) error
	DeleteProduct(c fiber.Ctx) error
}

type productsHandler struct {
	cfg             config.IConfig
	filesUsecase    filesUsecases.IFilesUsecase
	productsUsecase productsUsecases.IProductsUsecase
}

func ProductsHandler(cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase, productsUsecase productsUsecases.IProductsUsecase) IProductsHandler {
	return &productsHandler{
		cfg:             cfg,
		filesUsecase:    filesUsecase,
		productsUsecase: productsUsecase,
	}
}

// @Summary Find One Product
// @Description Find One Product
// @Tags Products
// @Accept  json
// @Produce  json
// @Param product_id path string true "Product ID"
// @Success 200 {object} products.Product
// @Router /products/{product_id} [get]
func (h *productsHandler) FindOneProduct(c fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	product, err := h.productsUsecase.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(findOneProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}

// @Summary Find Products
// @Description Find Products
// @Tags Products
// @Accept  json
// @Produce  json
// @Param id query string false "Id"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(10)
// @Param order_by query string false "Order By field" default(id)
// @Param sort_by query string false "Sort By direction (asc or desc)" default(desc)
// @Param search query string false "Search by title | description"
// @Success 200 {object} entities.PaginateRes
// @Router /products [get]
func (h *productsHandler) FindProduct(c fiber.Ctx) error {
	req := &products.ProductFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.Bind().Query(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(findProductErr),
			err.Error(),
		).Res()
	}

	if req.Page < 0 {
		req.Page = 1
	}

	if req.Limit < 5 {
		req.Limit = 5
	}

	if req.OrderBy == "" {
		req.OrderBy = "title"
	}

	if req.SortBy == "" {
		req.SortBy = "ASC"
	}

	products := h.productsUsecase.FindProduct(req)

	return entities.NewResponse(c).Success(fiber.StatusOK, products).Res()
}

// @Summary Add Product
// @Description Add Product
// @Tags Products
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body products.Product true "Product Request"
// @Success 200 {array} products.Product
// @Router /products [post]
func (h *productsHandler) AddProduct(c fiber.Ctx) error {
	req := &products.Product{
		Category: &appinfo.Category{},
		Images:   make([]*entities.Image, 0),
	}

	if err := c.Bind().JSON(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(insertProductErr),
			err.Error(),
		).Res()
	}

	if req.Category.Id <= 0 {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(insertProductErr),
			"category id is required",
		).Res()
	}

	product, err := h.productsUsecase.AddProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(insertProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, product).Res()
}

// @Summary Update Product
// @Description Update Product
// @Tags Products
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body products.Product true "Product Request"
// @Success 200 {array} products.Product
// @Router /products [patch]
func (h *productsHandler) UpdateProduct(c fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	req := &products.Product{
		Category: &appinfo.Category{},
		Images:   make([]*entities.Image, 0),
	}

	if err := c.Bind().JSON(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(updateProductErr),
			err.Error(),
		).Res()
	}
	req.Id = productId

	product, err := h.productsUsecase.UpdateProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(updateProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}

// @Summary Delete Product
// @Description Delete Product
// @Tags Products
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param product_id query string false "Product ID"
// @Success 200 {array} nil
// @Router /products/{product_id} [delete]
func (h *productsHandler) DeleteProduct(c fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	product, err := h.productsUsecase.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}

	deleteFileReq := make([]*files.DeleteFileReq, 0)
	for _, img := range product.Images {
		deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
			Destination: fmt.Sprintf("images/products/%s", img.Filename),
		})
	}

	if err := h.filesUsecase.DeleteFile(deleteFileReq); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}

	if err := h.productsUsecase.DeleteProduct(productId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}
	
	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}
