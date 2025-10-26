package controllers

import (
	"1001api/bookarena/pkg"
	"1001api/bookarena/views"

	"github.com/gofiber/fiber/v2"
)

type WebController struct{}

func NewWebController() *WebController {
	return &WebController{}
}

func (c *WebController) LoginPage(ctx *fiber.Ctx) error {
	return pkg.Render(ctx, views.LoginPage())
}

func (c *WebController) RegisterPage(ctx *fiber.Ctx) error {
	return pkg.Render(ctx, views.RegisterPage())
}
