package controllers

import (
	db "1001api/bookarena/internal/database/generated"
	"1001api/bookarena/internal/modules/payments"
	"1001api/bookarena/pkg"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/cast"
)

type PaymentController struct {
	paymentSvc payments.PaymentsService
	validate   *validator.Validate
}

func NewPaymentController(s payments.PaymentsService, v *validator.Validate) *PaymentController {
	return &PaymentController{
		paymentSvc: s,
		validate:   v,
	}
}

func (c *PaymentController) GetMasterPayments(ctx *fiber.Ctx) error {
	startTime := time.Now()

	limit := ctx.QueryInt("limit", 10)
	page := ctx.QueryInt("page", 1)

	result, err := c.paymentSvc.ListPaymentsForAdmin(limit, page)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.JSON(fiber.Map{
		"data": result,
		"meta": fiber.Map{
			"duration": time.Since(startTime).String(),
		},
	})
}

func (c *PaymentController) GetPaymentsByUser(ctx *fiber.Ctx) error {
	startTime := time.Now()

	userID := cast.ToString(ctx.Locals("user_id"))
	uid, err := uuid.Parse(userID)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"error": "unauthorized",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	limit := ctx.QueryInt("limit", 10)
	page := ctx.QueryInt("page", 1)

	result, err := c.paymentSvc.ListPaymentsForUser(uid, limit, page)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.JSON(fiber.Map{
		"data": result,
		"meta": fiber.Map{
			"duration": time.Since(startTime).String(),
		},
	})
}

func (c *PaymentController) GetPaymentByID(ctx *fiber.Ctx) error {
	startTime := time.Now()

	id := ctx.Params("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"error": "unauthorized",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	result, err := c.paymentSvc.GetPaymentByID(uid)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.JSON(fiber.Map{
		"data": result,
		"meta": fiber.Map{
			"duration": time.Since(startTime).String(),
		},
	})
}

func (c *PaymentController) DeletePayment(ctx *fiber.Ctx) error {
	startTime := time.Now()

	id := ctx.Params("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"error": "unauthorized",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	result, err := c.paymentSvc.DeletePayment(uid)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.JSON(fiber.Map{
		"data": result,
		"meta": fiber.Map{
			"duration": time.Since(startTime).String(),
		},
	})
}

func (c *PaymentController) PayDummyPayment(ctx *fiber.Ctx) error {
	startTime := time.Now()

	var req payments.PayInvoiceRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	if err := pkg.ValidateInput(req, c.validate); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err,
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	// Check if payment exists
	payment, err := c.paymentSvc.GetPaymentByID(req.InvoiceID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"error": "payment/invoice not found",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	// If payment is not pending, return error
	if payment.Status != string(pkg.PaymentStatusPending) {
		return ctx.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"error": "invalid payment status",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	// Update payment status
	if _, err := c.paymentSvc.UpdatePaymentStatus(db.UpdatePaymentStatusParams{
		ID:     req.InvoiceID,
		Status: string(pkg.PaymentStatusPaid),
	}); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.JSON(fiber.Map{
		"data": fiber.Map{
			"id":      payment.ID,
			"message": "Payment paid successfully",
		},
		"meta": fiber.Map{
			"duration": time.Since(startTime).String(),
		},
	})
}
