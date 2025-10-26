package controllers

import (
	"1001api/bookarena/internal/modules/fields"
	"1001api/bookarena/pkg"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type FieldController struct {
	fieldService fields.FieldService
	validator    *validator.Validate
}

func NewFieldController(s fields.FieldService, v *validator.Validate) *FieldController {
	return &FieldController{
		fieldService: s,
		validator:    v,
	}
}

func (c *FieldController) CreateField(ctx *fiber.Ctx) error {
	startTime := time.Now()

	var req fields.CreateFieldRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Error().Err(err).Msg("Failed to parse create field request")
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

	id, err := c.fieldService.CreateField(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create field")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "failed to create field",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"data": id,
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}

func (c *FieldController) GetMasterFields(ctx *fiber.Ctx) error {
	startTime := time.Now()

	limit := ctx.QueryInt("limit", 10)
	page := ctx.QueryInt("page", 1)

	fields, err := c.fieldService.ListFields(limit, page)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list fields")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "failed to list fields",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"data": fields,
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}

func (c *FieldController) GetFieldByID(ctx *fiber.Ctx) error {
	startTime := time.Now()

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse id")
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid id format",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	field, err := c.fieldService.GetFieldByID(id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return ctx.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"error": "field not found",
					"meta": fiber.Map{
						"duration": time.Since(startTime).String(),
					},
				},
			)
		}

		log.Error().Err(err).Msg("Failed to get field")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "failed to get field",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"data": field,
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}

func (c *FieldController) UpdateField(ctx *fiber.Ctx) error {
	startTime := time.Now()

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse id")
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid uuid format",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	var req fields.UpdateFieldRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Error().Err(err).Msg("Failed to parse update field request")
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

	updatedID, err := c.fieldService.UpdateField(id, req)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return ctx.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"error": "field not found",
					"meta": fiber.Map{
						"duration": time.Since(startTime).String(),
					},
				},
			)
		}

		log.Error().Err(err).Msg("Failed to update field")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "failed to update field",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"data": updatedID,
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}

func (c *FieldController) DeleteField(ctx *fiber.Ctx) error {
	startTime := time.Now()

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse id")
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid id format",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	deletedID, err := c.fieldService.DeleteField(id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return ctx.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"error": "field not found",
					"meta": fiber.Map{
						"duration": time.Since(startTime).String(),
					},
				},
			)
		}

		log.Error().Err(err).Msg("Failed to delete field")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "failed to delete field",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"data": deletedID,
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}
