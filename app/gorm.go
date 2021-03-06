package app

import (
	"errors"
	"github.com/arbrix/go-test/interfaces"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

type AppOrm struct {
	driver *gorm.DB
}

//Connect init gorm ORM.
func (orm *AppOrm) Connect(conf interfaces.Config) error {
	var db gorm.DB
	dbUri, err := conf.Get("DatabaseUri")
	if dbUri, ok := dbUri.(string); err == nil && ok {
		log.Println(dbUri)
		db, err = gorm.Open("mysql", dbUri)
	}
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	db.DB()

	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// Disable table name's pluralization
	db.SingularTable(true)
	if env, err := conf.Get("env"); err == nil && env == "dev" {
		db.LogMode(true)
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
