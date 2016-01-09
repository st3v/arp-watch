package observer

import (
	"encoding/json"
	"io/ioutil"
)

type MetronConfig struct {
	Origin   string `json:"origin"`
	Endpoint string `json:"endpoint"`
}

type Config struct {
	Metron    MetronConfig      `json:"metron,omitempty"`
	Frequency string            `json:"frequency"`
	Filters   []string          `json:"filters"`
	Aliases   map[string]string `json:"aliases"`
}

func (c *Config) Load(path string) error {
	return loadJSON(path, c)
}

func loadJSON(path string, dest interface{}) error {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(contents, dest); err != nil {
		return err
	}

	return nil
}
