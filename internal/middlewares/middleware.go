package middlewares

import (
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/contracts"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/jwt"
)

type Middleware struct {
	jwt         jwt.CustomJwtInterface
	sessionRepo contracts.SessionRepository
}

func NewMiddleware(
	jwt jwt.CustomJwtInterface,
	sessionRepo contracts.SessionRepository,
) *Middleware {
	return &Middleware{
		jwt:         jwt,
		sessionRepo: sessionRepo,
	}
}
