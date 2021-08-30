package services

import (
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

// Configuration ->
var Configuration *Config

// LogsConf ->
type LogsConf struct {
	Debug string `yaml:"debug"`
}

// DBConf ->
type DBConf struct {
	DBName          string        `yaml:"db_name"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Address         string        `yaml:"address"`
	Port            int           `yaml:"port"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// Config ->
type Config struct {
	Runas    string   `yaml:"-"`
	Port     string   `yaml:"port"`
	IP       string   `yaml:"ip"`
	StripeSK string   `yaml:"stripe_sk"`
	Logs     LogsConf `yaml:"logs"`
	DB       DBConf   `yaml:"db"`
}

var (
	// RunModeDev is a value for `runas` flag when service is running in development mode
	RunModeDev = "dev"
	// RunModeTest is a value for `runas` flag when service is running in test mode
	RunModeTest = "test"
	// RunModeProduction is a value for `runas` flag when service is running in production mode
	RunModeProduction = "prod"
)

// Configs ->
type Configs map[string]Config

// NewConfigurer ->
func NewConfigurer(runas, path string) *Config {
	fmt.Println("Using configuration for ", "'"+runas+"'", " deployment")

	data, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println("error opening configuration", err.Error())
	}

	var cs Configs

	err = yaml.Unmarshal(data, &cs)

	if err != nil {
		fmt.Println("Error unmarshaling ", runas, " configuration", err.Error())
	}

	c := cs[runas]
	c.Runas = runas
	Configuration = &c

	return &c
}
