package utils

import "github.com/google/uuid"

func GetUUID() string {
	newUUID := uuid.New()
	return newUUID.String()
}
