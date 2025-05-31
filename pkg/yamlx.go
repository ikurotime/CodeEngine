package pkg

import (
	"os"

	"gopkg.in/yaml.v2"
)

func ReadFile(filename string, structure interface{}) error {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlFile, structure)
}
