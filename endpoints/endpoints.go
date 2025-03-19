package endpoints

import (
	"log"
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Data struct {
	Endpoints map[string]map[string]interface{} `yaml:"endpoints"`
	Consumers []string                          `yaml:"consumers"`
}

var OperatingSystem string
var VaultUrl string
var ConfigData Data

func ReadConfig() {
	log.Println("Starting config file read task")
	status, _ := os.Stat("config.yml")
	if status == nil {
		log.Fatal(`Config file "config.yml" is missing`)
	}
	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatal("Failed due to ", err)
	}
	if ext := yaml.Unmarshal(data, &ConfigData); ext != nil {
		log.Fatal("Failed to unmarshal the config data")
	}

	OperatingSystem = runtime.GOOS
	if OperatingSystem == "windows" {
		VaultUrl = ConfigData.Endpoints["vault"]["addr1"].(string)
	}
	VaultUrl = ConfigData.Endpoints["vault"]["addr2"].(string)
}
