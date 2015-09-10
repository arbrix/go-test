package app

import (
	"errors"
	"github.com/arbrix/go-test/common"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type AppOrm struct {
	driver *gorm.DB
}

//Connect init gorm ORM.
func (orm *AppOrm) Connect(conf common.Config) error {
	var err error
	var db gorm.DB
	if dbUri, err := conf.Get("DatabaseUri"); err == nil {
		db, err = gorm.Open("mysql", dbUri)
	}
	if err != nil {
		return err
	}
	db.DB()

	// Then you could invoke `*sql.DB`'s functions with it
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// Disable table name's pluralization
	db.SingularTable(true)
	if env, err := conf.Get("env"); err == nil && env == "dev" {
		db.LogMode(true)
	}
	if err != nil {
		return err
	}
	orm.driver = &db
	return nil
}

func (orm *AppOrm) Create(model interface{}) error {
	if orm.driver.Create(model).Error != nil {
		return errors.New("The model is not created!")
	}
	return nil
}

func (orm *AppOrm) Find(models interface{}, vars interface{}) error {
	if orm.driver.Find(models, vars).RecordNotFound() {
		return errors.New("There is no one model found!")
	}
	return nil
}

func (orm *AppOrm) First(model interface{}, vars interface{}) error {
	if orm.driver.First(model, vars).RecordNotFound() {
		return errors.New("The model is not found!")
	}
	return nil
}

func (orm *AppOrm) Update(model interface{}, fieldSet map[string]interface{}) error {
	if orm.driver.Model(model).Update(fieldSet).Error != nil {
		return errors.New("The model is not updated!")
	}
	return nil
}

func (orm *AppOrm) Save(model interface{}) error {
	if orm.driver.Save(model).Error != nil {
		return errors.New("The model is not saved correctly!")
	}
	return nil
}

func (orm *AppOrm) Delete(model interface{}) error {
	if orm.driver.Delete(model).Error != nil {
		return errors.New("The model is not deleted correctly!")
	}
	return nil
}

func (orm *AppOrm) Close() error {
	return nil
}

func (orm *AppOrm) GetDriver() *gorm.DB {
	return orm.driver
}

type TestOrm struct {
}

//for in tests
func GetOrm4Test() (grm *AppOrm, err error) {
	grm = &AppOrm{}
	cnf := &TestConfig{}
	err = grm.Connect(cnf)
	return
}

func (orm *TestOrm) IsConnected() bool {
	return true
}

func (orm *TestOrm) Connect(cong common.Config) error {
	return nil
}

func (orm *TestOrm) Create(interface{}) error {
	return nil
}

func (orm *TestOrm) Find(interface{}, interface{}) error {
	return nil
}

func (orm *TestOrm) First(interface{}, interface{}) error {
	return nil
}

func (orm *TestOrm) Update(interface{}, map[string]interface{}) error {
	return nil
}

func (orm *TestOrm) Save(interface{}) error {
	return nil
}

func (orm *TestOrm) Delete(interface{}) error {
	return nil
}
