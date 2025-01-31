package api

import (
	"errors"

	"github.com/android-sms-gateway/ca/internal/csr"
	"github.com/android-sms-gateway/ca/pkg/core/handler"
	"github.com/android-sms-gateway/client-go/ca"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type csrHandler struct {
	handler.Base

	csrSvc *csr.Service
}

func (c *csrHandler) submit(ctx *fiber.Ctx) error {
	req := ca.PostCSRRequest{}
	if err := c.BodyParserValidator(ctx, &req); err != nil {
		return err
	}

	res, err := c.csrSvc.Create(ctx.Context(), csr.NewCSR(req.Content, req.Metadata))
	if err != nil {
		return err
	}

	return ctx.JSON(ca.PostCSRResponse{
		RequestID:   res.ID(),
		Status:      res.Status(),
		Message:     res.Status().Description(),
		Certificate: res.Certificate(),
	})
}

func (c *csrHandler) status(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	res, err := c.csrSvc.Get(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(ca.GetCSRStatusResponse{
		RequestID:   res.ID(),
		Status:      res.Status(),
		Message:     res.Status().Description(),
		Certificate: res.Certificate(),
	})
}

func (c *csrHandler) Register(router fiber.Router) {
	// router.Use(limiter.New(1, time.Minute))

	router.Use(c.handleError)

	router.Post("", c.submit)
	router.Get(":id", c.status)
}

func (c *csrHandler) handleError(ctx *fiber.Ctx) error {
	err := ctx.Next()

	if err == nil {
		return err
	}

	if errors.Is(err, csr.ErrCSRNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if errors.Is(err, csr.ErrCSRAlreadyExists) {
		return fiber.NewError(fiber.StatusConflict, err.Error())
	}

	if errors.Is(err, csr.ErrCSRInvalid) {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return err
}

func newCSR(csrSvc *csr.Service, v *validator.Validate, l *zap.Logger) *csrHandler {
	if csrSvc == nil {
		panic("service is required")
	}
	if v == nil {
		panic("validator is required")
	}
	if l == nil {
		panic("logger is required")
	}

	return &csrHandler{
		Base: handler.Base{
			Validator: v,
			Logger:    l,
		},

		csrSvc: csrSvc,
	}
}
