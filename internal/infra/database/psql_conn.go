package database

import (
	"fmt"
	"time"

	// pgx driver for postgres
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ahargunyllib/freepass-be-bcc-2025/internal/infra/env"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/log"

	"github.com/jmoiron/sqlx"
)

func NewPgsqlConn() *sqlx.DB {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable ",
		env.AppEnv.DBHost,
		env.AppEnv.DBPort,
		env.AppEnv.DBUser,
		env.AppEnv.DBPass,
		env.AppEnv.DBName,
	)

	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		log.Panic(log.LogInfo{
			"error": err.Error(),
		}, "[DB][NewPgsqlConn] failed to connect to database")
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}
