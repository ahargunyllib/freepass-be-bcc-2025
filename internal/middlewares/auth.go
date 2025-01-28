package middlewares

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/ahargunyllib/freepass-be-bcc-2025/domain"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (m *Middleware) RequireAuth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		header := ctx.Get("Authorization")
		if header == "" {
			return domain.ErrNoBearerToken
		}

		headerSlice := strings.Split(header, " ")
		if len(headerSlice) != 2 && headerSlice[0] != "Bearer" {
			return domain.ErrInvalidBearerToken
		}

		token := headerSlice[1]
		var claims jwt.Claims
		err := m.jwt.Decode(token, &claims)
		if err != nil {
			return domain.ErrInvalidBearerToken
		}

		notBefore, err := claims.GetNotBefore()
		if err != nil {
			return domain.ErrInvalidBearerToken
		}

		if notBefore.After(time.Now()) {
			return domain.ErrBearerTokenNotActive
		}

		expirationTime, err := claims.GetExpirationTime()
		if err != nil {
			return domain.ErrInvalidBearerToken
		}

		if expirationTime.Before(time.Now()) {
			return domain.ErrExpiredBearerToken
		}

		ctx.Locals("claims", claims)

		return ctx.Next()
	}
}

func (m *Middleware) AuthorizationSessionProposal() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claims, ok := ctx.Locals("claims").(jwt.Claims)
		if !ok {
			return domain.ErrNoBearerToken
		}

		if claims.Role == 2 { // event coordinator can access all session proposal
			return ctx.Next()
		}

		id := ctx.Params("id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			return domain.ErrSessionNotFound
		}

		session, err := m.sessionRepo.FindByID(ctx.Context(), uuid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return domain.ErrSessionNotFound
			}

			return err
		}

		if session.ProposerID != claims.UserID {
			return domain.ErrCantAccessResource
		}

		return ctx.Next()
	}
}

// If array is empty, it means all roles are accepted
func (m *Middleware) RequirePermission(
	acceptedRoles []int16,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claims, ok := ctx.Locals("claims").(jwt.Claims)
		if !ok {
			return domain.ErrNoBearerToken
		}

		if !m.isAccepted(claims.Role, acceptedRoles) {
			return domain.ErrCantAccessResource
		}

		return ctx.Next()
	}
}

func (m *Middleware) isAccepted(value int16, acceptedValues []int16) bool {
	if len(acceptedValues) == 0 {
		return true
	}
	for _, acceptedValue := range acceptedValues {
		if value == acceptedValue {
			return true
		}
	}
	return false
}
