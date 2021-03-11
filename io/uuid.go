package io

import "github.com/google/uuid"

func getUUID() string {
	return uuid.NewString()
}
