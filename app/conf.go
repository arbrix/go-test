package app

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
)

type AppConfig struct {
	basePath string
	options  map[string]interface{}
}

func NewAppConfig() *AppConfig {
	c := &AppConfig{}
	c.Load()
	return c
}

func (conf *AppConfig) Load() error {
	var env, basePath string
	flag.StringVar(&env, "env", "dev", "define environment: dev, prod, test (place config file *.json with the same name in ./config folder)")
	flag.StringVar(&basePath, "conf", "./config", "config path: ./config folder")
	flag.Parse()

	dir, err := filepath.Abs(basePath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	conf.basePath = dir
	log.Println(dir)
	conf.options = make(map[string]interface{})
	confPathSet := []string{conf.basePath + "/base.json"}
	confPathSet = append(confPathSet, conf.basePath+"/"+env+".json")
	for _, path := range confPathSet {
		src, err := conf.parsJson(path)
		if err != nil {
			return err
		}
		conf.mergeOpt(*src, true)
	}
	return nil
}

func (conf *AppConfig) Get(key string) (interface{}, error) {
	if val, ok := conf.options[key]; ok {
		return val, nil
	}
	return nil, errors.New("Options with key: " + key + " not exists! In " + conf.basePath)
}

func (conf *AppConfig) GetAll() *map[string]interface{} {
	return &conf.options
}

func (conf *AppConfig) parsJson(path string) (*map[string]interface{}, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, errors.New("config path not valid: " + path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	c := make(map[string]interface{})
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	return &c, err
}

func (conf *AppConfig) mergeOpt(src map[string]interface{}, rewrite bool) {
	for key, val := range src {
		if _, ok := conf.options[key]; ok == true && rewrite == false {
			continue //key already exists in option list
		}
		conf.options[key] = val
	}
}
