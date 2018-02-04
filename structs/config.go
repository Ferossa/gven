package structs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

const configVersion = "1"

type Config struct {
	Version      string                `json:"version"`
	Targets      map[string]Target     `json:"targets"`
	Repositories map[string]Repository `json:"repos"`
}

func (c *Config) Init(dir string) error {
	c.Version = configVersion
	c.Targets = make(map[string]Target)
	c.Repositories = make(map[string]Repository)

	return c.Save(dir)
}

func (c *Config) Save(dir string) error {
	configJson, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("error preparing config: %v", err)
	}

	configFile := filepath.Join(dir, "gven.json")
	err = ioutil.WriteFile(configFile, configJson, 0644)
	if err != nil {
		return fmt.Errorf("error saving config: %v", err)
	}

	return nil
}

type Target struct {
	Output       string            `json:"output"`
	Dependencies map[string]string `json:"deps"`
	Development  map[string]string `json:"dev"`
}

type Repository struct {
	Url  string `json:"url"`
	Type string `json:"type,omitempty"`
}
