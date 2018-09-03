package utils

import (
	"fmt"

	gouuid "github.com/satori/go.uuid"
)

func ValidateOrGenerateUUID(uuid string) string {
	_, err := gouuid.FromString(uuid)
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		uuid = GenerateUUID()
	}
	return uuid
}

func GenerateUUID() string {
	u := gouuid.Must(gouuid.NewV4())
	uuid := u.String()
	return uuid
}
