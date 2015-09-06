package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func makeConfig(env string) (*AppConfig, error) {
	conf := &AppConfig{basePath: "/Users/arbrix/dev/goproj/src/github.com/arbrix/go-test/config/"}
	err := conf.Load(env)
	return conf, err
}

func TestLoad(t *testing.T) {
	_, err := makeConfig("dev")
	if err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	conf, err := makeConfig("dev")
	val, err := conf.Get("api-url")
	assert.Equal(nil, err, "key 'api-url' must be in base.json")
	assert.Equal("/api/v1", val, "first verstion of API")
}
