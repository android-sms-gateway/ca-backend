package api

import (
	"errors"

	"github.com/android-sms-gateway/ca/internal/csr"
	"github.com/android-sms-gateway/client-go/ca"
	"github.com/android-sms-gateway/core/handler"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type csrHandler struct {
	handler.Base

	csrSvc *csr.Service
}

//	@Summary	Submit CSR
//	@Tags		CSR
//	@Accept		json
//	@Produce	json
//	@Param		request	body		ca.PostCSRRequest	true	"Request"
//	@Success	202		{object}	ca.PostCSRResponse
//	@Failure	400		{object}	http.ErrorResponse
//	@Failure	500		{object}	http.ErrorResponse
//	@Router		/csr [post]
//
// Submit CSR
func (c *csrHandler) submit(ctx *fiber.Ctx) error {
	req := ca.PostCSRRequest{}
	if err := c.BodyParserValidator(ctx, &req); err != nil {
		return err
	}

	res, err := c.csrSvc.Create(ctx.Context(), csr.NewCSR(req.Type, req.Content, req.Metadata))
	if err != nil {
		return err
	}

	return ctx.
		Status(fiber.StatusAccepted).
		JSON(csrStatusToResponse(res))
}

//	@Summary	Get CSR Status
//	@Tags		CSR
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Request ID"
//	@Success	200	{object}	ca.GetCSRStatusResponse
//	@Failure	400	{object}	http.ErrorResponse
//	@Failure	404	{object}	http.ErrorResponse
//	@Failure	500	{object}	http.ErrorResponse
//	@Router		/csr/{id} [get]
//
// Get CSR Status
func (c *csrHandler) status(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	res, err := c.csrSvc.Get(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(csrStatusToResponse(res))
}

func (c *csrHandler) Register(router fiber.Router) {
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
