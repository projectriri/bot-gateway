package utils

import (
	"fmt"
	gouuid "github.com/satori/go.uuid"
)

func ValidateOrGenerateUUID(uuid string) string {
	u, err := gouuid.FromString(uuid)
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		u = gouuid.Must(gouuid.NewV4())
		uuid = u.String()
	}
	return uuid
}