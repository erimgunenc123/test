package api_config

import (
	"fmt"
	"genericAPI/api/environment"
	"genericAPI/api/os_specific_constants"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var Config *config

func InitConfig() {
	wd, _ := os.Getwd()
	var configFileName string
	if environment.IsTestEnvironment() {
		configFileName = "test.yml"
	} else {
		configFileName = "prod.yml"
	}
	configFilePath := fmt.Sprintf("%s%sgenericAPI%sconfig%s%s", wd, os_specific_constants.PATH_SEPERATOR, os_specific_constants.PATH_SEPERATOR, os_specific_constants.PATH_SEPERATOR, configFileName)
	configBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		panic("Failed reading config file!")
	}

	err = yaml.Unmarshal(configBytes, &Config)
	if err != nil {
		panic("Failed unmarshalling config file!")
	}
	log.Print("Successfully initialized config")
}
