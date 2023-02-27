package config

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	jsonFile uint8 = iota
	yamlFile
	xmlFile
)

func ParseConfigFile(path string, config any) error {
	var t uint8

	switch filepath.Ext(path) {
	case ".json":
		t = jsonFile
	case ".yaml":
		t = yamlFile
	case ".xml":
		t = xmlFile
	}

	return parseFileConfig(path, t, config)
}

func parseFileConfig(path string, t uint8, config any) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch t {
	case jsonFile:
		return json.NewDecoder(file).Decode(config)
	case yamlFile:
		return yaml.NewDecoder(file).Decode(config)
	case xmlFile:
		return xml.NewDecoder(file).Decode(config)
	default:
		return errors.New("unsupported file type")
	}
}
