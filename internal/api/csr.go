package api

import (
	"github.com/android-sms-gateway/ca/pkg/client"
	"github.com/android-sms-gateway/ca/pkg/core/handler"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type csr struct {
	handler.Base
}

func (c *csr) submit(ctx *fiber.Ctx) error {
	req := client.PostCSRRequest{}
	if err := c.BodyParserValidator(ctx, &req); err != nil {
		return err
	}

	return ctx.JSON(client.PostCSRResponse{
		RequestID:   "123",
		Status:      client.CSRStatusApproved,
		Message:     client.CSRStatusDescriptionApproved,
		Certificate: "some cert",
	})
}

func (c *csr) status(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	return ctx.JSON(client.PostCSRResponse{
		RequestID:   id,
		Status:      client.CSRStatusApproved,
		Message:     client.CSRStatusDescriptionApproved,
		Certificate: "some cert",
	})
}

func (c *csr) Register(router fiber.Router) {
	// router.Use(limiter.New(1, time.Minute))

	router.Post("", c.submit)
	router.Get(":id", c.status)
}

func newCSR(v *validator.Validate, l *zap.Logger) *csr {
	return &csr{
		handler.Base{
			Validator: v,
			Logger:    l,
		},
	}
}
