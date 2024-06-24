package config

import (
	"log"

	_ "embed"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

//go:embed defaultConfig.yaml
var defaultConfigRaw []byte

func ReadConfig(config *koanf.Koanf, path string) {
	if err := config.Load(rawbytes.Provider(defaultConfigRaw), yaml.Parser()); err != nil {
		log.Fatal("Could not load default config\n")
	}

	config.Load(file.Provider(path), yaml.Parser())
}
