package middlewares

import (
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/log"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
)

func LoggerConfig() fiber.Handler {
	logger := log.GetLogger()
	config := fiberzerolog.Config{
		Logger:          logger,
		FieldsSnakeCase: true,
		Fields: []string{
			"referer",
			"ip",
			"host",
			"url",
			"ua",
			"latency",
			"status",
			"method",
			"error",
		},
		Messages: []string{
			"[LoggerMiddleware.LoggerConfig] Server error",
			"[LoggerMiddleware.LoggerConfig] Client error",
			"[LoggerMiddleware.LoggerConfig] Success",
		},
	}

	return fiberzerolog.New(config)
}
