package rest

import (
	"github.com/phathdt/service-context/component/validation"
	"mallbots/modules/user/application/dto"
	"mallbots/modules/user/domain/interfaces"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/phathdt/service-context/core"
)

type UserHandler struct {
	service interfaces.UserService
}

func NewUserHandler(service interfaces.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := validation.Validate(req); err != nil {
		panic(err)
	}

	token, err := h.service.Register(c.Context(), &req)
	if err != nil {
		panic(err)
	}

	return c.Status(http.StatusCreated).JSON(core.SimpleSuccessResponse(token))
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := validation.Validate(req); err != nil {
		panic(err)
	}

	token, err := h.service.Login(c.Context(), &req)
	if err != nil {
		panic(err)
	}

	return c.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(token))
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Context().UserValue("userId").(int32)

	profile, err := h.service.GetProfile(c.Context(), userID)
	if err != nil {
		panic(err)
	}

	return c.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(profile))
}
