package utils

import (
"encoding/json"
"log"
)

// PrettyPrint prints a struct in a pretty format
func PrettyPrint(v interface{}) {
b, err := json.MarshalIndent(v, "", "  ")
if err != nil {
log.Println("Error:", err)
}
log.Println(string(b))
}
