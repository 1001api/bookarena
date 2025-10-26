package controllers

import (
	"1001api/bookarena/internal/modules/bookings"
	"1001api/bookarena/pkg"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
)

type BookingController struct {
	bookingService bookings.BookingService
	validator      *validator.Validate
}

func NewBookingController(s bookings.BookingService, v *validator.Validate) *BookingController {
	return &BookingController{
		bookingService: s,
		validator:      v,
	}
}

func (c *BookingController) CreateBooking(ctx *fiber.Ctx) error {
	startTime := time.Now()

	// current user id
	id := cast.ToString(ctx.Locals("user_id"))
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

	var req bookings.CreateBookingRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Error().Err(err).Msg("failed to parse json body")
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid json body",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	if err := pkg.ValidateInput(req, c.validator); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err,
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	createdID, err := c.bookingService.CreateBooking(uid, req)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return ctx.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"error": "field not found",
					"meta": fiber.Map{
						"duration": time.Since(startTime).String(),
					},
				},
			)
		}

		if strings.Contains(err.Error(), "booking conflict") {
			return ctx.Status(fiber.StatusConflict).JSON(
				fiber.Map{
					"error": "Booking conflict with already existing booking, please choose another timeframe",
					"meta": fiber.Map{
						"duration": time.Since(startTime).String(),
					},
				},
			)
		}

		log.Error().Err(err).Msg("failed to create booking")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"data": createdID,
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}

func (c *BookingController) GetBookingByID(ctx *fiber.Ctx) error {
	startTime := time.Now()

	id := cast.ToString(ctx.Params("id"))
	uid, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid id",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	booking, err := c.bookingService.GetBookingByID(uid)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return ctx.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"error": "booking not found",
					"meta": fiber.Map{
						"duration": time.Since(startTime).String(),
					},
				},
			)
		}

		log.Error().Err(err).Msg("failed to get booking")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "failed to get booking",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"data": booking,
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}

func (c *BookingController) GetBookingsByUser(ctx *fiber.Ctx) error {
	startTime := time.Now()

	id := cast.ToString(ctx.Locals("user_id"))
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

	limit := ctx.QueryInt("limit", 10)
	page := ctx.QueryInt("page", 1)

	bookings, err := c.bookingService.ListBookingsByUser(uid, limit, page)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user bookings")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "failed to get user bookings",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"data": bookings,
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}

func (c *BookingController) GetBookingsByField(ctx *fiber.Ctx) error {
	startTime := time.Now()

	id := cast.ToString(ctx.Params("id"))
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

	limit := ctx.QueryInt("limit", 10)
	page := ctx.QueryInt("page", 1)

	bookings, err := c.bookingService.ListBookingsByField(uid, limit, page)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user bookings")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "failed to get user bookings",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"data": bookings,
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}

func (c *BookingController) GetMasterBookings(ctx *fiber.Ctx) error {
	startTime := time.Now()

	limit := ctx.QueryInt("limit", 10)
	page := ctx.QueryInt("page", 1)

	bookings, err := c.bookingService.ListBookings(limit, page)
	if err != nil {
		log.Error().Err(err).Msg("failed to get master bookings")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "failed to get master bookings",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"data": bookings,
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}

func (c *BookingController) DeleteBooking(ctx *fiber.Ctx) error {
	startTime := time.Now()

	id := cast.ToString(ctx.Params("id"))
	uid, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid booking id",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	bookingID, err := c.bookingService.DeleteBooking(uid)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return ctx.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"error": "booking not found",
					"meta": fiber.Map{
						"duration": time.Since(startTime).String(),
					},
				},
			)
		}

		log.Error().Err(err).Msg("failed to delete booking")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "failed to delete booking",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"data": bookingID,
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}
