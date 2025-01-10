package config

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

var Conf = &GocnConfig{}

type GocnConfig struct {
	Pool     map[string]any
	Template map[string]any
}

func init() {
	loadToml()
}

func loadToml() {
	configFile := flag.String("conf", "conf/app.toml", "app config file")
	flag.Parse()
	if _, err := os.Stat(*configFile); err != nil {
		fmt.Println("conf/app.toml file not load，because not exist")
		return
	}
	_, err := toml.DecodeFile(*configFile, Conf)
	if err != nil {
		fmt.Println("conf/app.toml decode fail check format")
		return
	}
}
