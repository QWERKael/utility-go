package config

import (
	"github.com/QWERKael/utility-go/io"
	"gopkg.in/yaml.v2"
)

func ParserFromByte(configText []byte, conf interface{}) error {
	err := yaml.Unmarshal(configText, conf)
	return err
}

func ParserFromPath(configPath string, conf interface{}) error {
	b, err := io.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = ParserFromByte(b, conf)
	return err
}
