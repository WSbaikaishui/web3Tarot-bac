package util

import "github.com/google/uuid"

func UuidV4() string {
	return uuid.New().String()
}
