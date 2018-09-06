package utils

import (
	gouuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

func ValidateOrGenerateUUID(uuid string) string {
	_, err := gouuid.FromString(uuid)
	if err != nil {
		log.Warn("[uuid] %s", err)
		uuid = GenerateUUID()
	}
	return uuid
}

func GenerateUUID() string {
	u := gouuid.Must(gouuid.NewV4())
	uuid := u.String()
	return uuid
}
