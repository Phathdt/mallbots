package middleware

import (
	"github.com/gofiber/fiber/v2"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/core"
	"mallbots/plugins/tokenprovider"
	common2 "mallbots/shared/common"
	"strings"

	"github.com/pkg/errors"
)

func ExtractTokenFromHeaderString(headers []string) (string, error) {
	if len(headers) == 0 {
		return "", errors.New("missing token")
	}
	//"Authorization" : "Bearer {token}"

	parts := strings.Split(headers[0], " ")

	if len(parts) == 0 {
		return "", errors.New("missing token")
	}

	if parts[0] != "Bearer" || len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "", errors.New("wrong authen header")
	}

	return parts[1], nil
}

func RequiredAuth(sc sctx.ServiceContext) fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()
		token, err := ExtractTokenFromHeaderString(headers["Authorization"])

		if err != nil {
			panic(core.ErrUnauthorized.WithError(err.Error()))
		}

		tokenProvider := sc.MustGet(common2.KeyJwt).(tokenprovider.Provider)

		payload, err := tokenProvider.Validate(token)
		if err != nil {
			panic(core.ErrUnauthorized.WithError(err.Error()))
		}

		c.Context().SetUserValue("userId", payload.GetUserId())
		return c.Next()
	}
}
