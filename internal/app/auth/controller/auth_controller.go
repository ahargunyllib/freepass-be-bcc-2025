package controller

import (
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain"
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/contracts"
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/dto"
	"github.com/ahargunyllib/freepass-be-bcc-2025/internal/middlewares"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/helpers/http/response"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

type authController struct {
	service contracts.AuthService
}

func InitAuthController(router fiber.Router, service contracts.AuthService, middleware *middlewares.Middleware) {
	controller := authController{
		service: service,
	}

	authRouter := router.Group("/auth")

	authRouter.Post("/login", controller.Login)
	authRouter.Post("/register", controller.Register)
	authRouter.Get("/session", middleware.RequireAuth(), controller.GetSession)
}

func (a *authController) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	err := a.service.Register(c.Context(), req)
	if err != nil {
		return err
	}

	return response.SendResponse(c, fiber.StatusCreated, nil)
}

func (a *authController) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	res, err := a.service.Login(c.Context(), req)
	if err != nil {
		return err
	}

	return response.SendResponse(c, fiber.StatusOK, res)
}

func (a *authController) GetSession(c *fiber.Ctx) error {
	claims, ok := c.Locals("claims").(jwt.Claims)
	if !ok {
		return domain.ErrExpiredBearerToken
	}

	query := dto.SessionQuery{
		UserID: claims.UserID,
	}

	res, err := a.service.GetSession(c.Context(), query)
	if err != nil {
		return err
	}

	return response.SendResponse(c, fiber.StatusOK, res)
}
