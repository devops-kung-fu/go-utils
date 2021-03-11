package io

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getUUID(t *testing.T) {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	uuid := getUUID()
	assert.Regexp(t, r, uuid, "Not a valid UUID")
}
