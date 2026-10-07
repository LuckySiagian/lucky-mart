package utils

import (
	"github.com/google/uuid"
)

func GenerateUUID() uuid.UUID {
	return uuid.New()
}

func ParseUUID(str string) (uuid.UUID, error) {
	return uuid.Parse(str)
}

func IsValidUUID(str string) bool {
	_, err := uuid.Parse(str)
	return err == nil
}
