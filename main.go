package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"

	"github.com/arbrix/go-test/config"
	"github.com/arbrix/go-test/web"
)

func getConfig(path string) (config.Config, error) {
	jsonPath := path
	conf := config.Config{}

	if _, err := os.Stat(jsonPath); err != nil {
		return conf, errors.New("config path not valid")
	}

	file, err := os.Open(jsonPath)
	if err != nil {
		return conf, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	log.Println(conf.String())
	return conf, err
}

func main() {
	var confPath string

	flag.StringVar(&confPath, "config", "conf.json", "path to config file")
	flag.Parse()

	var err error
	config.JsonConfig, err = getConfig(confPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	svc := web.Service{}

	svc.Run(config.JsonConfig)
}
