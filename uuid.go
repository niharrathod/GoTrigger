package main

import (
	"log"

	uuid "github.com/nu7hatch/gouuid"
)

// GetNewUniqueID generates new unique id
func GetNewUniqueID() string {
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
		panic("failed to generate UUID")
	}
	return uuid.String()
}
