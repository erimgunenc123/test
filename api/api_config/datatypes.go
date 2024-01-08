package api_config

import "fmt"

type config struct {
	DB  dbConfig  `yaml:"db"`
	App appConfig `yaml:"app"`
}

type appConfig struct {
	Port    string        `yaml:"port"`
	Logging loggingConfig `yaml:"logging"`
}

type loggingConfig struct {
	LogFilepath string `yaml:"log_filepath"`
}

type dbConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"database"`
}

func (db *dbConfig) GetConnectionStr() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db.Username, db.Password, db.Host, db.Port, db.DatabaseName)
}
