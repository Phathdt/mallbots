package rest

import (
	"mallbots/modules/product/application/dto"
	"mallbots/modules/product/domain/interfaces"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/phathdt/service-context/core"
)

type ProductHandler struct {
	service interfaces.ProductService
}

func NewProductHandler(service interfaces.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	type reqParam struct {
		dto.ProductListRequest
		core.Paging
	}

	var rp reqParam

	if err := c.QueryParser(&rp); err != nil {
		panic(err)
	}

	rp.Paging.Process()

	products, err := h.service.GetProducts(c.Context(), &rp.ProductListRequest, &rp.Paging)
	if err != nil {
		panic(err)
	}

	return c.Status(http.StatusOK).JSON(core.ResponseWithPaging(products, &rp.ProductListRequest, &rp.Paging))
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		panic(err)
	}

	product, err := h.service.GetProduct(c.Context(), int32(id))
	if err != nil {
		panic(err)
	}

	return c.Status(http.StatusOK).JSON(core.SimpleSuccessResponse(product))
}
