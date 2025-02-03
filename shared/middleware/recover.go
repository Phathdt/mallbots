package middleware

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
	"net/http"
)

type CanGetStatusCode interface {
	StatusCode() int
}

func Recover(serviceCtx sctx.ServiceContext) fiber.Handler {
	// Return new handler
	return func(ctx *fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				if appErr, ok := err.(CanGetStatusCode); ok {
					ctx.Status(appErr.StatusCode()).JSON(appErr)
				} else {
					if fieldErrors, ok := err.(validator.ValidationErrors); ok {
						message := getMessageError(fieldErrors)

						ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
							"code":    http.StatusBadRequest,
							"message": message,
						})

					} else {
						ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
							"code":    http.StatusInternalServerError,
							"status":  "internal server error",
							"message": "something went wrong, please try again or contact supporters",
						})
					}
				}

				serviceCtx.Logger("service").Errorf("%+v \n", err)
			}
		}()

		// Return err if existed, else move to next handler
		return ctx.Next()
	}
}

func getMessageError(fieldErrors []validator.FieldError) string {
	fieldError := fieldErrors[0]

	//TODO: add more tag
	switch fieldError.Tag() {
	case "required":
		return fmt.Sprintf("%s is a required field", fieldError.Field())
	case "max":
		return fmt.Sprintf("%s must be a maximum of %s in length", fieldError.Field(), fieldError.Param())
	case "min":
		return fmt.Sprintf("%s must be a minimum of %s in length", fieldError.Field(), fieldError.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", fieldError.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid Email", fieldError.Field())
	case "oneof":
		return fmt.Sprintf("%s must be one of enums %s", fieldError.Field(), fieldError.Param())
	case "hourtime":
		return fmt.Sprintf("%s must be between 00:00 and 23:59", fieldError.Field())
	case "requirethenmust":
		return fmt.Sprintf("leng %s must be %s", fieldError.Field(), fieldError.Param())
	case "gtcsfield":
		return fmt.Sprintf("%s must be greater than %s", fieldError.Field(), fieldError.Param())
	default:
		return fmt.Sprintf("something wrong on %s; %s", fieldError.Field(), fieldError.Tag())
	}
}
