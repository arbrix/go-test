package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testConfig struct {
}

func (cnf *testConfig) Load(env string) error {
	return nil
}

func (cnf *testConfig) Get(key string) (interface{}, error) {
	return "test:test@/test?charset=utf8&parseTime=True&loc=Local", nil
}

func (cnf *testConfig) GetAll() *map[string]interface{} {
	return nil
}

func TestConnect(t *testing.T) {
	_, err := getOrm()
	assert.Equal(t, nil, err, "error while connection to DB")
}

type Test struct {
	Title  string
	Status int
}

func TestCreate(t *testing.T) {
	grm, err := getOrm()
	test := &Test{}
	db := grm.GetDriver()
	err = db.CreateTable(&Test{}).Error
	assert.Equal(t, nil, err, "error while careate table in DB")
	err = grm.Create(test)
	assert.Equal(t, nil, err, "error while write to DB")
	//err = db.DropTable(&Test{}).Error
	//assert.Equal(t, nil, err, "error while droped table in DB")

}

func getOrm() (grm *AppOrm, err error) {
	grm = &AppOrm{}
	cnf := &testConfig{}
	err = grm.Connect(cnf)
	return
}
