package app

import (
	"encoding/json"
	"errors"
	"os"
)

//Config interface that describe methods for works with it
type Config interface {
	Load(env string) error
	Get(key string) (interface{}, error)
	GetAll() *map[string]interface{}
}

type AppConfig struct {
	basePath string
	options  map[string]interface{}
}

func (conf *AppConfig) Load(env string) error {
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

type TestConfig struct {
}

func (cnf *TestConfig) Load(env string) error {
	return nil
}

func (cnf *TestConfig) Get(key string) (interface{}, error) {
	switch key {
	case "env":
		return "dev", nil
	case "DatabaseUri":
		return "test:test@/test?charset=utf8&parseTime=True&loc=Local", nil
	case "SecretKey":
		return "sdkfaSWerjPDFEjiRwErfjOSDFIj39024@#4()urr2,nasroiu3@#$I23Sf0(Ur23ks0f9@#rjSf0W#rjl23kng0-)I#l23n", nil
	default:
		return nil, errors.New("key: " + key + " not defined!")
	}
}

func (cnf *TestConfig) GetAll() *map[string]interface{} {
	return nil
}
