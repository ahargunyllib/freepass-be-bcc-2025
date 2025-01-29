package server

import (
	authController "github.com/ahargunyllib/freepass-be-bcc-2025/internal/app/auth/controller"
	authRepo "github.com/ahargunyllib/freepass-be-bcc-2025/internal/app/auth/repository"
	authSvc "github.com/ahargunyllib/freepass-be-bcc-2025/internal/app/auth/service"
	sessionController "github.com/ahargunyllib/freepass-be-bcc-2025/internal/app/session/controller"
	sessionRepo "github.com/ahargunyllib/freepass-be-bcc-2025/internal/app/session/repository"
	sessionSvc "github.com/ahargunyllib/freepass-be-bcc-2025/internal/app/session/service"
	userController "github.com/ahargunyllib/freepass-be-bcc-2025/internal/app/user/controller"
	userRepo "github.com/ahargunyllib/freepass-be-bcc-2025/internal/app/user/repository"
	userSvc "github.com/ahargunyllib/freepass-be-bcc-2025/internal/app/user/service"
	"github.com/ahargunyllib/freepass-be-bcc-2025/internal/middlewares"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/bcrypt"
	errorhandler "github.com/ahargunyllib/freepass-be-bcc-2025/pkg/helpers/http/error_handler"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/helpers/http/response"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/jwt"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/log"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/uuid"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/validator"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type HTTPServer interface {
	Start(part string)
	MountMiddlewares()
	MountRoutes(db *sqlx.DB)
	GetApp() *fiber.App
}

type httpServer struct {
	app *fiber.App
}

func NewHTTPServer() HTTPServer {
	config := fiber.Config{
		CaseSensitive: true,
		AppName:       "BCC Conference (Freepass BE)",
		ServerHeader:  "BCC Conference (Freepass BE)",
		JSONEncoder:   sonic.Marshal,
		JSONDecoder:   sonic.Unmarshal,
		ErrorHandler:  errorhandler.ErrorHandler,
	}

	app := fiber.New(config)

	return &httpServer{
		app: app,
	}
}

func (s *httpServer) GetApp() *fiber.App {
	return s.app
}

func (s *httpServer) Start(port string) {
	if port[0] != ':' {
		port = ":" + port
	}

	err := s.app.Listen(port)

	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[SERVER][Start] failed to start server")
	}
}

func (s *httpServer) MountMiddlewares() {
	s.app.Use(middlewares.LoggerConfig())
	s.app.Use(middlewares.Helmet())
	s.app.Use(middlewares.Compress())
	s.app.Use(middlewares.Cors())
	s.app.Use(middlewares.RecoverConfig())
}

func (s *httpServer) MountRoutes(db *sqlx.DB) {
	bcrypt := bcrypt.Bcrypt
	uuid := uuid.UUID
	validator := validator.Validator
	jwt := jwt.Jwt

	s.app.Get("/", func(c *fiber.Ctx) error {
		return response.SendResponse(c, fiber.StatusOK, "Freepass BE BCC 2025")
	})

	api := s.app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return response.SendResponse(c, fiber.StatusOK, "Freepass BE BCC 2025")
	})

	userRepository := userRepo.NewUserRepository(db)
	authRepository := authRepo.NewAuthRepository(db)
	sessionRepository := sessionRepo.NewSessionRepository(db)

	middleware := middlewares.NewMiddleware(jwt, sessionRepository)

	userService := userSvc.NewUserService(userRepository, validator, uuid, bcrypt)
	authService := authSvc.NewAuthService(authRepository, validator, uuid, bcrypt, jwt)
	sessionService := sessionSvc.NewSessionService(sessionRepository, validator, uuid)

	userController.InitUserController(v1, userService, middleware)
	authController.InitAuthController(v1, authService, middleware)
	sessionController.InitSessionController(v1, sessionService, middleware)

	s.app.Use(func(c *fiber.Ctx) error {
		return c.SendFile("./web/not-found.html")
	})
}
