package main

import (
	"github.com/ahargunyllib/freepass-be-bcc-2025/internal/infra/database"
	"github.com/ahargunyllib/freepass-be-bcc-2025/internal/infra/env"
	"github.com/ahargunyllib/freepass-be-bcc-2025/internal/infra/server"
)

func main() {
	server := server.NewHTTPServer()
	psqlDB := database.NewPgsqlConn()
	defer psqlDB.Close()

	server.MountMiddlewares()
	server.MountRoutes(psqlDB)
	server.Start(env.AppEnv.AppPort)
}
