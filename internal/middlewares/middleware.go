package middlewares

import "github.com/ahargunyllib/freepass-be-bcc-2025/pkg/jwt"

type Middleware struct {
	jwt jwt.CustomJwtInterface
}

func NewMiddleware(
	jwt jwt.CustomJwtInterface,
) *Middleware {
	return &Middleware{
		jwt: jwt,
	}
}
