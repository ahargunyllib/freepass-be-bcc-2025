package uuid

import (
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/log"
	"github.com/google/uuid"
)

type CustomUUIDInterface interface {
	NewV7() (uuid.UUID, error)
}

type CustomUUIDStruct struct{}

var UUID = getUUID()

func getUUID() CustomUUIDInterface {
	return &CustomUUIDStruct{}
}

func (u *CustomUUIDStruct) NewV7() (uuid.UUID, error) {
	uuid, err := uuid.NewV7()

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[UUID][New] failed to create uuid v7")

		return uuid, err
	}

	return uuid, err
}
