package api_config

import "fmt"

type config struct {
	DB      dbConfig      `yaml:"db"`
	App     appConfig     `yaml:"app"`
	Btcturk btcturkConfig `yaml:"btcturk"`
}

type btcturkConfig struct {
	PublicKey  string `yaml:"public_key"`
	PrivateKey string `yaml:"private_key"`
}

type appConfig struct {
	Port    string        `yaml:"port"`
	Logging loggingConfig `yaml:"logging"`
	Secret  string        `yaml:"secret"`
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
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", db.Username, db.Password, db.Host, db.Port, db.DatabaseName)
}
