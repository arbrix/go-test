package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"

	"github.com/arbrix/go-test/web"
)

func getConfig(path string) (web.Config, error) {
	jsonPath := path
	config := web.Config{}

	if _, err := os.Stat(jsonPath); err != nil {
		return config, errors.New("config path not valid")
	}

	file, err := os.Open(jsonPath)
	if err != nil {
		return config, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	log.Println(config.String())
	return config, err
}

func main() {
	var confPath string

	flag.StringVar(&confPath, "config", "conf.json", "path to config file")
	flag.Parse()

	cfg, err := getConfig(confPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	svc := web.TaskService{}

	if err = svc.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
