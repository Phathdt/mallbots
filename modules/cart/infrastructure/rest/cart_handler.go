package rest

import (
	"mallbots/modules/cart/application/dto"
	"mallbots/modules/cart/domain/interfaces"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/phathdt/service-context/component/validation"
	"github.com/phathdt/service-context/core"
)

type CartHandler struct {
	service interfaces.CartService
}

func NewCartHandler(service interfaces.CartService) *CartHandler {
	return &CartHandler{service: service}
}

func (h *CartHandler) AddItem(c *fiber.Ctx) error {
	var req dto.CartItemRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := validation.Validate(req); err != nil {
		panic(err)
	}

	userID := c.Context().UserValue("userId").(int32)

	item, err := h.service.AddItem(c.Context(), userID, &req)
	if err != nil {
		panic(err)
	}

	return c.Status(http.StatusCreated).JSON(core.SimpleSuccessResponse(item))
}

func (h *CartHandler) UpdateQuantity(c *fiber.Ctx) error {
	var req dto.CartItemRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := validation.Validate(req); err != nil {
		panic(err)
	}

	userID := c.Context().UserValue("userId").(int32)

	item, err := h.service.UpdateQuantity(c.Context(), userID, &req)
	if err != nil {
		panic(err)
	}

	return c.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(item))
}

func (h *CartHandler) RemoveItem(c *fiber.Ctx) error {
	userID := c.Context().UserValue("userId").(int32)
	productID, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		panic(err)
	}

	if err := h.service.RemoveItem(c.Context(), userID, int32(productID)); err != nil {
		panic(err)
	}

	return c.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(true))
}

func (h *CartHandler) GetItems(c *fiber.Ctx) error {
	userID := c.Context().UserValue("userId").(int32)

	items, err := h.service.GetItems(c.Context(), userID)
	if err != nil {
		panic(err)
	}

	return c.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(items))
}
