package app

import (
	"encoding/json"
	"errors"
	"os"
)

//Config interface that describe methods for works with it
type Config interface {
	Init(path string) error
	Get(key string) (interface{}, error)
	GetAll() *map[string]interface{}
}

type AppConfig struct {
	basePath string
	options  map[string]interface{}
}

func (conf *AppConfig) Init(env string) error {
	conf.options = make(map[string]interface{})
	confPathSet := []string{conf.basePath + "base.json"}
	confPathSet = append(confPathSet, conf.basePath+env+".json")
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
	return nil, errors.New("Options with key: " + key + " not exists!")
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
