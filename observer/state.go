package observer

import (
	"encoding/json"
	"io/ioutil"
)

type State map[string]string

func (s *State) Load(path string) error {
	return loadJSON(path, s)
}

func (s *State) Write(path string) error {
	contents := []byte("{}")

	if len(*s) > 0 {
		var err error
		if contents, err = json.MarshalIndent(s, "", "  "); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(path, contents, 0644)
}
