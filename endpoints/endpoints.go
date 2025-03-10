package endpoints

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Data struct {
	Endpoints map[string]map[string]interface{} `yaml:"endpoints"`
	Consumers []string                          `yaml:"consumers"`
}

var ConfigData Data

func ReadConfig() {
	log.Println("Starting config file read task")
	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatal("Failed due to ", err)
	}
	if ext := yaml.Unmarshal(data, &ConfigData); ext != nil {
		log.Fatal("Failed to unmarshal the config data")
	}
}
