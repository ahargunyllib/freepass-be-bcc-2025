package controller

import (
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain"
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/contracts"
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/dto"
	"github.com/ahargunyllib/freepass-be-bcc-2025/internal/middlewares"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/helpers/http/response"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type sessionController struct {
	service contracts.SessionService
}

func InitSessionController(router fiber.Router, service contracts.SessionService, middleware *middlewares.Middleware) {
	controller := sessionController{service: service}

	sessionRouter := router.Group("/sessions")
	sessionRouter.Get("/", middleware.RequireAuth(),
		middleware.RequirePermission([]int16{1, 2}), // user, event coordinator
		controller.GetSessions)
	sessionRouter.Get("/:id", controller.GetSession)
	sessionRouter.Post(
		"/",
		middleware.RequireAuth(),
		middleware.RequirePermission([]int16{1}), // user
		controller.CreateSession,
	)
	sessionRouter.Put(
		"/:id",
		middleware.RequireAuth(),
		middleware.RequirePermission([]int16{1}), // user
		middleware.AuthorizationSessionProposal(),
		controller.UpdateSession,
	)
	sessionRouter.Delete(
		"/:id",
		middleware.RequireAuth(),
		middleware.RequirePermission([]int16{1}), // user
		middleware.AuthorizationSessionProposal(),
		controller.DeleteSession,
	)

	sessionRouter.Post(
		"/:id/accept",
		middleware.RequireAuth(),
		middleware.RequirePermission([]int16{2}), // admin
		controller.AcceptSession,
	)
	sessionRouter.Post(
		"/:id/reject",
		middleware.RequireAuth(),
		middleware.RequirePermission([]int16{2}), // admin
		controller.RejectSession,
	)
	sessionRouter.Post(
		"/:id/cancel",
		middleware.RequireAuth(),
		middleware.RequirePermission([]int16{2}), // admin
		controller.CancelSession,
	)


	sessionRouter.Post(
		"/:session_id/register",
		middleware.RequireAuth(),
		middleware.RequirePermission([]int16{1}), // user
		controller.RegisterSession,
	)
	sessionRouter.Post(
		"/:sesion_id/unregister",
		middleware.RequireAuth(),
		middleware.RequirePermission([]int16{1}), // user
		controller.UnregisterSession,
	)
	sessionRouter.Post(
		"/:sesion_id/review",
		middleware.RequireAuth(),
		middleware.RequirePermission([]int16{1}), // user
		controller.ReviewSession,
	)
	sessionRouter.Post(
		"/:sesion_id/review/:user_id/remove",
		middleware.RequireAuth(),
		middleware.RequirePermission([]int16{2}), // admin
		controller.RemoveReview,
	)
}

/*
All user can get all sessions (not his) THAT is ACCEPTED
User can get all his sessions
Admin can get all sessions
*/
func (c *sessionController) GetSessions(ctx *fiber.Ctx) error {
	var query dto.GetSessionsQuery
	if err := ctx.QueryParser(&query); err != nil {
		return err
	}

	// get user's proposals
	if query.ProposerID != uuid.Nil {
		claims, ok := ctx.Locals("claims").(jwt.Claims)
		if !ok {
			return response.SendResponse(ctx, fiber.StatusUnauthorized, nil)
		}

		// only allow user to get their own proposals or admin to get all proposals
		if claims.Role != 2 && query.ProposerID != claims.UserID {
			return domain.ErrCantAccessResource
		}
	}

	// get not accepted proposals
	if query.Status != 2 { // not accepted
		claims, ok := ctx.Locals("claims").(jwt.Claims)
		if !ok {
			return response.SendResponse(ctx, fiber.StatusUnauthorized, nil)
		}

		// user want to get all someone not accepted proposals
		if claims.Role == 1 && query.ProposerID != uuid.Nil && query.ProposerID != claims.UserID {
			return domain.ErrCantAccessResource
		}

		// specify the proposer id if user want to get all not accepted proposals but not giving the proposer id
		if claims.Role == 1 && query.ProposerID == uuid.Nil {
			query.ProposerID = claims.UserID
		}
	}

	sessions, err := c.service.GetSessions(ctx.Context(), query)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, sessions)
}

func (c *sessionController) GetSession(ctx *fiber.Ctx) error {
	var query dto.GetSessionEventQuery
	if err := ctx.ParamsParser(&query); err != nil {
		return err
	}

	session, err := c.service.GetSession(ctx.Context(), query)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, session)
}

func (c *sessionController) CreateSession(ctx *fiber.Ctx) error {
	var req dto.CreateSessionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return response.SendResponse(ctx, fiber.StatusUnauthorized, nil)
	}

	req.ProposerID = claims.UserID

	err := c.service.CreateSession(ctx.Context(), req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusCreated, nil)
}

func (c *sessionController) UpdateSession(ctx *fiber.Ctx) error {
	var req dto.UpdateSessionRequest

	if err := ctx.ParamsParser(&req); err != nil {
		return err
	}

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	err := c.service.UpdateSession(ctx.Context(), req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, nil)
}

func (c *sessionController) DeleteSession(ctx *fiber.Ctx) error {
	var query dto.DeleteSessionQuery
	if err := ctx.ParamsParser(&query); err != nil {
		return err
	}

	err := c.service.DeleteSession(ctx.Context(), query)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, nil)
}

func (c *sessionController) AcceptSession(ctx *fiber.Ctx) error {
	var query dto.AcceptSessionQuery
	if err := ctx.ParamsParser(&query); err != nil {
		return err
	}

	var req dto.AcceptSessionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	err := c.service.AcceptSession(ctx.Context(), query, req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, nil)
}

func (c *sessionController) RejectSession(ctx *fiber.Ctx) error {
	var query dto.RejectSessionQuery
	if err := ctx.ParamsParser(&query); err != nil {
		return err
	}

	var req dto.RejectSessionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	err := c.service.RejectSession(ctx.Context(), query, req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, nil)
}

func (c *sessionController) RegisterSession(ctx *fiber.Ctx) error {
	var query dto.RegisterSessionQuery
	if err := ctx.ParamsParser(&query); err != nil {
		return err
	}

	var req dto.RegisterSessionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return response.SendResponse(ctx, fiber.StatusUnauthorized, nil)
	}

	req.UserID = claims.UserID

	err := c.service.RegisterSession(ctx.Context(), query, req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, nil)
}

func (c *sessionController) UnregisterSession(ctx *fiber.Ctx) error {
	var query dto.UnregisterSessionQuery
	if err := ctx.ParamsParser(&query); err != nil {
		return err
	}

	var req dto.UnregisterSessionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return response.SendResponse(ctx, fiber.StatusUnauthorized, nil)
	}

	req.UserID = claims.UserID

	err := c.service.UnregisterSession(ctx.Context(), query, req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, nil)
}

func (c *sessionController) ReviewSession(ctx *fiber.Ctx) error {
	var query dto.ReviewSessionQuery
	if err := ctx.ParamsParser(&query); err != nil {
		return err
	}

	var req dto.ReviewSessionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	claims, ok := ctx.Locals("claims").(jwt.Claims)
	if !ok {
		return response.SendResponse(ctx, fiber.StatusUnauthorized, nil)
	}

	req.UserID = claims.UserID

	err := c.service.ReviewSession(ctx.Context(), query, req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, nil)
}

func (c *sessionController) RemoveReview(ctx *fiber.Ctx) error {
	var query dto.DeleteReviewSessionQuery
	if err := ctx.ParamsParser(&query); err != nil {
		return err
	}

	var req dto.DeleteReviewSessionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	err := c.service.DeleteReviewSession(ctx.Context(), query, req)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, nil)
}

func (c *sessionController) CancelSession(ctx *fiber.Ctx) error {
	var query dto.CancelSessionQuery
	if err := ctx.ParamsParser(&query); err != nil {
		return err
	}

	err := c.service.CancelSession(ctx.Context(), query)
	if err != nil {
		return err
	}

	return response.SendResponse(ctx, fiber.StatusOK, nil)
}
