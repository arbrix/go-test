package app

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Orm interface {
	Connect(cong Config) error
	Create(interface{}) error
	Find(interface{}, int64) error
	Update(interface{}, map[string]interface{}) error
	Save(interface{}) error
	Delete(interface{}) error
}

type AppOrm struct {
	db gorm.DB
}

// Init init gorm ORM.
func (orm *AppOrm) Connect(conf Config) error {
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
	orm.db = db
	return nil
}

func (orm *AppOrm) Create(model interface{}) error {
	if orm.db.Create(model).Error != nil {
		return errors.New("The model is not created!")
	}
	return nil
}

func (orm *AppOrm) Find(model interface{}, id int64) error {
	if orm.db.First(model, id).RecordNotFound() {
		return errors.New("The model is not found!")
	}
	return nil
}

func (orm *AppOrm) Update(model interface{}, fieldSet map[string]interface{}) error {
	if orm.db.Model(model).Update(fieldSet).Error != nil {
		return errors.New("The model is not updated!")
	}
	return nil
}

func (orm *AppOrm) Save(model interface{}) error {
	if orm.db.Save(model).Error != nil {
		return errors.New("The model is not saved correctly!")
	}
	return nil
}

func (orm *AppOrm) Delete(model interface{}) error {
	if orm.db.Delete(model).Error != nil {
		return errors.New("The model is not deleted correctly!")
	}
	return nil
}

func (orm *AppOrm) Close() error {
	return nil
}

func (orm *AppOrm) GetDriver() *gorm.DB {
	return &orm.db
}
