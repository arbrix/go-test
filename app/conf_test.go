package app

import (
	"github.com/stretchr/testify/assert"
	"log"
	"path/filepath"
	"testing"
)

type TConf struct {
	AppConfig
}

func (c *TConf) Load() error {
	//change this vars for tests
	env, basePath := "dev", "../config" //../ - for go-test/app package tests, ./ - for go-test general testing
	dir, err := filepath.Abs(basePath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	c.basePath = dir
	c.options = make(map[string]interface{})
	confPathSet := []string{c.basePath + "/base.json"}
	confPathSet = append(confPathSet, c.basePath+"/"+env+".json")
	for _, path := range confPathSet {
		src, err := c.parsJson(path)
		if err != nil {
			return err
		}
		c.mergeOpt(*src, true)
	}
	return nil
}

func makeConfig() (*TConf, error) {
	conf := &TConf{}
	err := conf.Load()
	return conf, err
}

func TestLoad(t *testing.T) {
	_, err := makeConfig()
	if err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	conf, err := makeConfig()
	val, err := conf.Get("api-url")
	assert.Equal(nil, err, "key 'api-url' must be in base.json")
	assert.Equal("/api/v1", val, "first verstion of API")
}
