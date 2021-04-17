package configuration

import (
	"context"
	"fmt"
	"net"

	"github.com/sethvargo/go-envconfig"
)

type Configuration struct {
	Address    string `env:"ADDRESS,required"`
	ListenPort int    `env:"PORT,default=4560"`
	ListenAddr string `env:"LISTEN_ADDR,default=0.0.0.0"`
	DataDir    string `env:"DATA_DIR,default=/data"`
}

var config Configuration

func Config() Configuration {
	return config
}

func LoadConfigFromEnv() error {
	ctx := context.Background()

	if err := envconfig.Process(ctx, &config); err != nil {
		return err
	}

	return checkConfiguration()
}

func checkConfiguration() error {
	_, _, err := net.SplitHostPort(fmt.Sprintf("%s:1", config.ListenAddr))
	if err != nil {
		return fmt.Errorf("Invalid listen address: %s", config.ListenAddr)
	}
	return nil
}
