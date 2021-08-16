package config

import (
	"fmt"
	"log"
	"path"
	"testing"
)

type Conf struct {
	Host   string `yaml:"host"`
	Listen int    `yaml:"listen"`
}

func TestParserFromByte(t *testing.T) {
	var configText = `
host: 127.0.0.1
listen: 4001
`
	conf := Conf{}
	err := ParserFromByte([]byte(configText), &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%#v\n\n", conf)
}

func TestParserFromPath(t *testing.T) {
	configPath := path.Join("..", "testDir", "testConfig.yml")
	conf := Conf{}
	err := ParserFromPath(configPath, &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%#v\n\n", conf)
}
