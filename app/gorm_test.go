package app

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testConfig struct {
}

func (cnf *testConfig) Load(env string) error {
	return nil
}

func (cnf *testConfig) Get(key string) (interface{}, error) {
	switch key {
	case "env":
		return "dev", nil
	case "DatabaseUri":
		return "test:test@/test?charset=utf8&parseTime=True&loc=Local", nil
	}
	return nil, errors.New("key: " + key + " not defined!")
}

func (cnf *testConfig) GetAll() *map[string]interface{} {
	return nil
}

func TestConnect(t *testing.T) {
	_, err := getOrm()
	assert.Equal(t, nil, err, "error while connection to DB")
}

type Test struct {
	gorm.Model
	Title  string
	Status int
}

func (t Test) TableName() string {
	return "test"
}

func TestCreate(t *testing.T) {
	grm, test, _ := testSetup(t)
	db := grm.GetDriver()
	assert.Equal(t, false, db.NewRecord(test), "table must contain this model already")
	testShutdown(t, db)
}

func TestFind(t *testing.T) {
	grm, test, err := testSetup(t)
	searchTest := Test{}
	err = grm.Find(&searchTest, int64(test.ID))
	assert.Equal(t, nil, err, "error while reading from DB")
	assert.Equal(t, test.Title, searchTest.Title, "error while comparing saved and finded element from DB")
	testShutdown(t, grm.GetDriver())
}

func TestUpdate(t *testing.T) {
	grm, test, err := testSetup(t)
	err = grm.Update(test, map[string]interface{}{"status": 15})
	assert.Equal(t, nil, err, "error while updating record in DB")
	assert.Equal(t, 15, test.Status, "Status must be updated and qeal 15")
	testShutdown(t, grm.GetDriver())
}

func TestSave(t *testing.T) {
	grm, test, err := testSetup(t)
	test.Status = 15
	err = grm.Save(test)
	assert.Equal(t, nil, err, "error while saving record in DB")
	assert.Equal(t, 15, test.Status, "Status must be updated and qeal 15")
	testShutdown(t, grm.GetDriver())
}

func TestDelete(t *testing.T) {
	grm, test, err := testSetup(t)
	err = grm.Delete(test)
	assert.Equal(t, nil, err, "error while deleting the record from DB")
	db := grm.GetDriver()
	assert.Equal(t, false, db.NewRecord(test), "the record was deleted from DB and key mast be free")
	testShutdown(t, grm.GetDriver())
}

func getOrm() (grm *AppOrm, err error) {
	grm = &AppOrm{}
	cnf := &testConfig{}
	err = grm.Connect(cnf)
	return
}

func testSetup(t *testing.T) (grm *AppOrm, test *Test, err error) {
	grm, err = getOrm()
	db := grm.GetDriver()
	err = db.CreateTable(&Test{}).Error
	assert.Equal(t, nil, err, "error while careate table in DB")
	test = &Test{Title: "test", Status: 10}
	assert.Equal(t, true, db.NewRecord(test), "table must be empty, and primire kay is free")
	err = grm.Create(test)
	assert.Equal(t, nil, err, "error while write to DB task")
	return
}

func testShutdown(t *testing.T, db *gorm.DB) {
	err := db.DropTable(&Test{}).Error
	assert.Equal(t, nil, err, "error while droped table in DB")
}
