package main

import (
	"encoding/csv"
	"os"

	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/entity"
	"github.com/ahargunyllib/freepass-be-bcc-2025/internal/infra/database"
	"github.com/ahargunyllib/freepass-be-bcc-2025/internal/infra/env"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/bcrypt"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/log"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/uuid"
	"github.com/jmoiron/sqlx"
)

const SeedersFilePath = "data/seeders/"
const SeedersDevPath = SeedersFilePath + "dev/"
const SeedersProdPath = SeedersFilePath + "prod/"

func main() {
	psqlDB := database.NewPgsqlConn()
	defer psqlDB.Close()

	var path string
	if env.AppEnv.AppEnv == "production" {
		path = SeedersProdPath
	} else {
		path = SeedersDevPath
	}

	uuid := uuid.UUID
	bcrypt := bcrypt.Bcrypt

	seedAdmins(path, psqlDB, uuid, bcrypt)
}

func seedAdmins(path string, db *sqlx.DB, uuid uuid.CustomUUIDInterface, bcrypt bcrypt.CustomBcryptInterface) {
	path += "admins.csv"

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err,
		}, "[seed][seedUsers] Error opening file")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Error(log.LogInfo{
			"error": err,
		}, "[seed][seedUsers] Error reading file")
		return
	}

	for idx, record := range records {
		if idx == 0 { // skipping header
			continue
		}

		log.Info(log.LogInfo{
			"record": record,
		}, "[seed][seedUsers] Inserting record")

		id, err := uuid.NewV7()
		if err != nil {
			log.Error(log.LogInfo{
				"error": err,
			}, "[seed][seedUsers] Error generating UUID")
		}

		hashedPassword, err := bcrypt.Hash(record[2])
		if err != nil {
			log.Error(log.LogInfo{
				"error": err,
			}, "[seed][seedUsers] Error hashing password")
		}

		user := entity.User{
			ID:       id,
			Name:     record[0],
			Email:    record[1],
			Password: hashedPassword,
			Role:     3, // admin
		}

		_, err = db.NamedExec(
			`
			INSERT INTO users (id, name, email, role, password)
			VALUES (:id, :name, :email, :role, :password)
			`,
			user,
		)

		if err != nil {
			log.Error(log.LogInfo{
				"error": err,
			}, "[seed][seedUsers] Error inserting admins")
			return
		}
	}
}
