package util

import (
	"fmt"

	"github.com/google/uuid"
)

func NewV1UUID() string {

	id, err := uuid.NewUUID()

	if err != nil {
		fmt.Printf("%s\n", err)
		return ""
	}

	return id.String()
}

func NewV4UUID() string {
	id := uuid.New()
	return id.String()
}
