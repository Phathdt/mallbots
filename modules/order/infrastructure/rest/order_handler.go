package rest

import (
	"mallbots/modules/order/application/dto"
	"mallbots/modules/order/domain/interfaces"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/phathdt/service-context/component/validation"
	"github.com/phathdt/service-context/core"
)

type OrderHandler struct {
	service interfaces.OrderService
}

func NewOrderHandler(service interfaces.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var req dto.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := validation.Validate(req); err != nil {
		panic(err)
	}

	userID := c.Context().UserValue("userId").(int32)

	order, err := h.service.CreateOrder(c.Context(), userID, &req)
	if err != nil {
		panic(err)
	}

	return c.Status(http.StatusCreated).JSON(core.SimpleSuccessResponse(order))
}

func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		panic(core.ErrBadRequest.WithError(err.Error()))
	}

	order, err := h.service.GetOrder(c.Context(), int32(id))
	if err != nil {
		panic(err)
	}

	return c.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(order))
}

func (h *OrderHandler) GetUserOrders(c *fiber.Ctx) error {
	type reqParam struct {
		core.Paging
	}

	var rp reqParam
	if err := c.QueryParser(&rp); err != nil {
		panic(err)
	}

	rp.Paging.Process()

	userID := c.Context().UserValue("userId").(int32)

	orders, err := h.service.GetUserOrders(c.Context(), userID, &rp.Paging)
	if err != nil {
		panic(err)
	}

	return c.Status(http.StatusOK).JSON(core.ResponseWithPaging(orders, nil, &rp.Paging))
}
