package utils

import (
	"encoding/json"
	"log"
)

// PrettyPrint prints a struct in a pretty format
func PrettyPrint(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}
	log.Println(string(b))
}
