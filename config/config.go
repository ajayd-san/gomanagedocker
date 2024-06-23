package config

import (
	"log"
	"os"

	_ "embed"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

//go:embed defaultConfig.yaml
var defaultConfigRaw []byte

func ReadConfig(config *koanf.Koanf) {
	//read config file
	configPath, err := os.UserConfigDir()

	if err != nil {
		log.Println("$HOME could not be determined")
	}

	if err := config.Load(rawbytes.Provider(defaultConfigRaw), yaml.Parser()); err != nil {
		log.Fatal("Could not load default config\n")
	}

	config.Load(file.Provider(configPath+xdgPathTail), yaml.Parser())
}
