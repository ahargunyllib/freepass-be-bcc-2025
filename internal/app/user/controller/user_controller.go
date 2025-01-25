package controller

import (
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/contracts"
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/dto"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/helpers/http/response"
	"github.com/gofiber/fiber/v2"
)

type userController struct {
	userService contracts.UserService
}

func InitUserController(router fiber.Router, userService contracts.UserService) {
	controller := userController{
		userService: userService,
	}

	userRouter := router.Group("/users")

	userRouter.Get("/", controller.GetUsers)
	userRouter.Get("/:id", controller.GetUser)
	userRouter.Post("/", controller.CreateUser)
	userRouter.Delete("/:id", controller.DeleteUser)
}

func (u *userController) GetUsers(c *fiber.Ctx) error {
	var query dto.GetUsersQuery
	if err := c.QueryParser(&query); err != nil {
		return err
	}

	users, err := u.userService.GetUsers(c.Context(), query)
	if err != nil {
		return err
	}

	return response.SendResponse(c, fiber.StatusOK, users)
}

func (u *userController) GetUser(c *fiber.Ctx) error {
	var query dto.GetUserQuery
	if err := c.ParamsParser(&query); err != nil {
		return err
	}

	user, err := u.userService.GetUser(c.Context(), query)
	if err != nil {
		return err
	}

	return response.SendResponse(c, fiber.StatusOK, user)
}

func (u *userController) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	err := u.userService.CreateUser(c.Context(), req)
	if err != nil {
		return err
	}

	return response.SendResponse(c, fiber.StatusCreated, nil)
}

func (u *userController) DeleteUser(c *fiber.Ctx) error {
	var query dto.DeleteUserQuery
	if err := c.ParamsParser(&query); err != nil {
		return err
	}

	err := u.userService.DeleteUser(c.Context(), query)
	if err != nil {
		return err
	}

	return response.SendResponse(c, fiber.StatusOK, nil)
}
