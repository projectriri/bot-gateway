package utils

import (
	gouuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

func ValidateOrGenerateUUID(uuid string) string {
	if !ValidateUUID(uuid) {
		uuid = GenerateUUID()
	}
	return uuid
}

func GenerateUUID() string {
	u := gouuid.Must(gouuid.NewV4())
	uuid := u.String()
	return uuid
}

func ValidateUUID(uuid string) bool {
	_, err := gouuid.FromString(uuid)
	if err != nil {
		log.Warnf("[uuid] %s", err)
		return false
	}
	return true
}
