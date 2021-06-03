package util

import "github.com/google/uuid"

func getUUID() string {
	return uuid.NewString()
}
