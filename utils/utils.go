package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// WriteJsonToFile writes data to path+filename.json
func WriteJsonToFile(path string, filename string, data interface{}) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(fmt.Sprintf("%s/%s.json", path, filename), json, 0644)
}
