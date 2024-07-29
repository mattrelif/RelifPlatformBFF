package utils

import "github.com/google/uuid"

type UuidGenerator func() string

func GenerateUuid() string {
	return uuid.NewString()
}
