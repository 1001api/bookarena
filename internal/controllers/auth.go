package controllers

import (
	"1001api/bookarena/internal/modules/auth"
	"1001api/bookarena/pkg"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type AuthController struct {
	authService auth.AuthService
	validator   *validator.Validate
}

func NewAuthController(s auth.AuthService, v *validator.Validate) *AuthController {
	return &AuthController{
		authService: s,
		validator:   v,
	}
}

func (c *AuthController) RegisterUser(ctx *fiber.Ctx) error {
	startTime := time.Now()

	var req auth.RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
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

	// Set role to public/user since this method is called from public endpoint
	req.Role = string(pkg.RoleUser)

	createdID, err := c.authService.Register(req)
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
		"data": createdID.String(),
		"meta": fiber.Map{
			"duration": time.Since(startTime).String(),
		},
	})
}

func (c *AuthController) RegisterAdmin(ctx *fiber.Ctx) error {
	startTime := time.Now()

	var req auth.RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
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

	// Set role to admin since this method is called from admin endpoint
	req.Role = string(pkg.RoleAdmin)

	createdID, err := c.authService.Register(req)
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
		"data": createdID.String(),
		"meta": fiber.Map{
			"duration": time.Since(startTime).String(),
		},
	})
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	startTime := time.Now()

	var req auth.LoginRequest

	if err := ctx.BodyParser(&req); err != nil {
		log.Error().Err(err).Msg("Failed to parse login request")
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid json body",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	resp, err := c.authService.Login(req)
	if err != nil {
		if strings.Contains(err.Error(), "restricted access") {
			return ctx.Status(fiber.StatusUnauthorized).JSON(
				fiber.Map{
					"error": "restricted access",
					"meta": fiber.Map{
						"duration": time.Since(startTime).String(),
					},
				},
			)
		}

		return ctx.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"error": err.Error(),
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     pkg.AccessTokenName,
		Value:    resp.Token,
		Path:     "/",
		Domain:   viper.GetString("APP_DOMAIN"),
		HTTPOnly: true,
		Secure:   viper.GetBool("ENV_PROD"),
		SameSite: fiber.CookieSameSiteLaxMode,
		MaxAge:   60 * viper.GetInt("JWT_EXPIRY"), // 15 minutes
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     pkg.RefreshTokenName,
		Value:    resp.RefreshToken,
		Path:     "/",
		Domain:   viper.GetString("APP_DOMAIN"),
		HTTPOnly: true,
		Secure:   viper.GetBool("ENV_PROD"),
		SameSite: fiber.CookieSameSiteLaxMode,
		MaxAge:   60 * viper.GetInt("JWT_REFRESH_EXPIRY"), // 3 days
	})

	return ctx.JSON(
		fiber.Map{
			"data": fiber.Map{
				"token":         resp.Token,
				"refresh_token": resp.RefreshToken,
			},
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}

func (c *AuthController) Refresh(ctx *fiber.Ctx) error {
	startTime := time.Now()

	refreshToken := ctx.Cookies(pkg.RefreshTokenName)

	if refreshToken == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"error": "refresh_token is required",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	resp, err := c.authService.RefreshToken(auth.RefreshRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"error": err.Error(),
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     pkg.AccessTokenName,
		Value:    resp.Token,
		Path:     "/",
		Domain:   viper.GetString("APP_DOMAIN"),
		HTTPOnly: true,
		Secure:   viper.GetBool("ENV_PROD"),
		SameSite: fiber.CookieSameSiteLaxMode,
		MaxAge:   60 * viper.GetInt("JWT_EXPIRY"), // 15 minutes default
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     pkg.RefreshTokenName,
		Value:    resp.RefreshToken,
		Path:     "/",
		Domain:   viper.GetString("APP_DOMAIN"),
		HTTPOnly: true,
		Secure:   viper.GetBool("ENV_PROD"),
		SameSite: fiber.CookieSameSiteLaxMode,
		MaxAge:   60 * viper.GetInt("JWT_REFRESH_EXPIRY"), // 3 days default
	})

	return ctx.JSON(
		fiber.Map{
			"data": fiber.Map{
				"token":         resp.Token,
				"refresh_token": resp.RefreshToken,
			},
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}
