package ordersHandlers

import (
	"strings"
	"time"

	"github.com/IzePhanthakarn/kawaii-shop/config"
	"github.com/IzePhanthakarn/kawaii-shop/modules/entities"
	"github.com/IzePhanthakarn/kawaii-shop/modules/orders"
	"github.com/IzePhanthakarn/kawaii-shop/modules/orders/ordersUsecases"
	"github.com/gofiber/fiber/v3"
)

type ordersHandlersErrCode string

const (
	findOneOrderErr ordersHandlersErrCode = "orders-001"
	findOrderErr    ordersHandlersErrCode = "orders-002"
	insertOrderErr  ordersHandlersErrCode = "orders-003"
)

type IOrdersHandler interface {
	FindOneOrder(c fiber.Ctx) error
	FindOrder(c fiber.Ctx) error
	InsertOrder(c fiber.Ctx) error
}

type ordersHandlers struct {
	cfg          config.IConfig
	orderUsecase ordersUsecases.IOrdersUsecase
}

func OrdersHandlers(cfg config.IConfig, orderUsecase ordersUsecases.IOrdersUsecase) IOrdersHandler {
	return &ordersHandlers{
		cfg:          cfg,
		orderUsecase: orderUsecase,
	}
}

func (h *ordersHandlers) FindOneOrder(c fiber.Ctx) error {
	orderId := strings.Trim(c.Params("order_id"), " ")

	order, err := h.orderUsecase.FindOneOrder(orderId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(findOneOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}

func (h *ordersHandlers) FindOrder(c fiber.Ctx) error {
	req := &orders.OrderFilter{
		SortReq:       &entities.SortReq{},
		PaginationReq: &entities.PaginationReq{},
	}

	if err := c.Bind().Query(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(findOrderErr),
			err.Error(),
		).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 5 {
		req.Limit = 5
	}

	orderByMap := map[string]string{
		"id":         `"o"."id"`,
		"created_at": `"o"."created_at"`,
	}
	if orderByMap[req.SortBy] == "" {
		req.OrderBy = orderByMap["id"]
	}

	req.SortBy = strings.ToUpper(req.SortBy)
	sortMap := map[string]string{
		"ASC":  "ASC",
		"DESC": "DESC",
	}
	if sortMap[req.SortBy] != "" {
		req.SortBy = sortMap["DESC"]
	}

	if req.StartDate != "" {
		start, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(findOrderErr),
				"invalid start date",
			).Res()
		}
		req.StartDate = start.Format("2006-01-02")
	}

	if req.EndDate != "" {
		end, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(findOrderErr),
				"invalid end date",
			).Res()
		}
		req.EndDate = end.Format("2006-01-02")
	}

	orders := h.orderUsecase.FindOrder(req)

	return entities.NewResponse(c).Success(fiber.StatusOK, orders).Res()
}

func (h *ordersHandlers) InsertOrder(c fiber.Ctx) error {
	userId := strings.Trim(c.Locals("userId").(string), " ")

	req := &orders.Order{
		Products: make([]*orders.ProductsOrder, 0),
	}

	if err := c.Bind().JSON(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(insertOrderErr),
			err.Error(),
		).Res()
	}

	if len(req.Products) == 0 {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(insertOrderErr),
			"request body is empty",
		).Res()
	}

	if c.Locals("userRoleId").(int) != 2 {
		req.UserId = userId
	}

	req.Status = "waiting"
	req.TotalPaid = 0

	order, err := h.orderUsecase.InsertOrder(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(insertOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}
