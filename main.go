package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/arbrix/go-test/service"
)

func getConfig(path string) (service.Config, error) {
	jsonPath := path
	config := service.Config{}

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
	confPath := os.Args[1]
	cfg, err := getConfig(confPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	svc := service.TaskService{}

	if err = svc.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
