package configs

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ptonlix/netdog/pkg/file"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var config = new(Config)

type Config struct {
	Network struct {
		Node     string `toml:"node"`
		Cron     string `toml:"cron"`
		PingTest struct {
			Enabled      bool `toml:"enabled"`
			PingInterval int  `toml:"pinginterval"`
			Device       []struct {
				Name string `toml:"name"`
				Ip   string `toml:"ip"`
			} `toml:"device"`
		} `toml:"pingtest"`
		BindwidthTest struct {
			Enabled bool `toml:"enabled"`
			Device  []struct {
				Name string `toml:"name"`
				Ip   string `toml:"ip"`
			} `toml:"device"`
		} `toml:"bindwidthtest"`
	} `toml:"network"`
	Mail struct {
		Host string `toml:"host"`
		Port int    `toml:"port"`
		User string `toml:"user"`
		Pass string `toml:"pass"`
		To   string `toml:"to"`
	} `toml:"mail"`
	Data struct {
		Pingfile string `toml:"pingfile"`
	} `toml:"data"`
}

var (
	//go:embed netdog_configs.toml
	netdogConfigs []byte
)

func init() {
	var r io.Reader

	r = bytes.NewReader(netdogConfigs)

	viper.SetConfigType("toml")

	if err := viper.ReadConfig(r); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(config); err != nil {
		panic(err)
	}
	fmt.Println(config)
	viper.SetConfigName(ProjectName + "_configs")
	viper.AddConfigPath("./configs")

	configFile := "./configs/" + ProjectName + "_configs.toml"
	_, ok := file.IsExists(configFile)
	if !ok {
		if err := os.MkdirAll(filepath.Dir(configFile), 0766); err != nil {
			panic(err)
		}

		f, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := viper.WriteConfig(); err != nil {
			panic(err)
		}
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(config); err != nil {
			panic(err)
		}
	})
}

func Get() Config {
	return *config
}
