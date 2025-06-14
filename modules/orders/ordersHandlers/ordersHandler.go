package ordersHandlers

import (
	"strings"
	"time"

	"github.com/IzePhanthakarn/go-basic-shop/config"
	"github.com/IzePhanthakarn/go-basic-shop/modules/entities"
	"github.com/IzePhanthakarn/go-basic-shop/modules/orders"
	"github.com/IzePhanthakarn/go-basic-shop/modules/orders/ordersUsecases"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ordersHandlersErrCode string

const (
	findOneOrderErr ordersHandlersErrCode = "orders-001"
	findOrderErr    ordersHandlersErrCode = "orders-002"
	insertOrderErr  ordersHandlersErrCode = "orders-003"
	updateOrderErr  ordersHandlersErrCode = "orders-004"
)

type IOrdersHandler interface {
	FindOneOrder(c fiber.Ctx) error
	FindOrder(c fiber.Ctx) error
	InsertOrder(c fiber.Ctx) error
	UpdateOrder(c fiber.Ctx) error
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

// @Summary Find One Order
// @Description Find One Order
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param user_id path string true "User ID"
// @Param order_id path string true "Order ID"
// @Security BearerAuth
// @Success 200 {object} orders.Order
// @Router /orders/{user_id}/{order_id} [get]
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

// @Summary Find Orders
// @Description Find Orders
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(10)
// @Param order_by query string false "Order By field" default(id)
// @Param sort_by query string false "Sort By direction (asc or desc)" default(desc)
// @Param search query string false "Search by user_id | address | contact"
// @Param status query string false "Status"
// @Param start_date query string false "Start Date (YYYY-MM-DD)"
// @Param end_date query string false "End Date (YYYY-MM-DD)"
// @Security BearerAuth
// @Success 200 {object} entities.PaginateRes
// @Router /orders [get]
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

// @Summary Insert Order
// @Description Insert Order
// @Tags Orders
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body orders.OrderReq true "Order Request"
// @Success 200 {array} orders.Order
// @Router /orders [post]
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

// @Summary Update Order
// @Description Update Order
// @Tags Orders
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param user_id query string false "User ID"
// @Param order_id query string false "Order ID"
// @Param request body orders.OrderReq true "Order Request"
// @Success 200 {array} orders.Order
// @Router /orders/{user_id}/{order_id} [patch]
func (h *ordersHandlers) UpdateOrder(c fiber.Ctx) error {
	orderId := strings.Trim(c.Params("order_id"), " ")
	req := new(orders.Order)
	if err := c.Bind().JSON(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(updateOrderErr),
			err.Error(),
		).Res()
	}
	req.Id = orderId

	statusMap := map[string]string{
		"waiting":   "waiting",
		"shipping":  "shipping",
		"completed": "completed",
		"canceled":  "canceled",
	}

	if c.Locals("userRoleId").(int) != 2 {
		req.UserId = strings.Trim(c.Locals("userId").(string), " ")
	} else if strings.ToLower(req.Status) == statusMap["canceles"] {
		req.Status = statusMap["canceled"]
	}

	if req.TransferSlip != nil {
		if req.TransferSlip.Id == "" {
			req.TransferSlip.Id = uuid.NewString()
		}
		if req.TransferSlip.CreatedAt == "" {
			loc, err := time.LoadLocation("Asia/Bangkok")
			if err != nil {
				return entities.NewResponse(c).Error(
					fiber.StatusInternalServerError,
					string(updateOrderErr),
					err.Error(),
				).Res()
			}
			now := time.Now().In(loc)

			// YYYY-MM-DD HH:MM:SS
			// 2006-01-02 15:04:05
			req.TransferSlip.CreatedAt = now.Format("2006-01-02 15:04:05")
		}
	}

	order, err := h.orderUsecase.UpdateOrder(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(updateOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}
