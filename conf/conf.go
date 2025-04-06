package conf

import (
	"os"
	"sync"

	"github.com/google/wire"
	"gopkg.in/yaml.v3"
)

var (
	once       sync.Once
	config     *Config
	ProvideSet = wire.NewSet(New)
)

type Config struct {
	TCPConfig	`yaml:"tcp"`
}

type TCPConfig struct {
	Addr string `yaml:"addr"`
}

func readConfig() {
	// panic("to be completed")
	data, err := os.ReadFile("conf.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
}

func New() *Config {
	once.Do(readConfig)
	return config
}
