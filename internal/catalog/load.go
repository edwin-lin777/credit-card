package catalog

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// read one card file
func LoadCard(path string) (Card, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Card{}, fmt.Errorf("parsing %s: %w", path, err)
	}
	var c Card
	if err := yaml.Unmarshal(data, &c); err != nil {
		return Card{}, fmt.Errorf("parsing %s: %w", path, err)
	}
	return c, nil

}
