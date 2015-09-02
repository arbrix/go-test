package app

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Orm interface {
	Init(cong Config)
	Find(*gorm.Model, int64) error
	Save(*gorm.Model) error
	Delete(*gorm.Model) error
}

type AppOrm struct {
	db gorm.DB
}

// Init init gorm ORM.
func (orm *AppOrm) Init(conf AppConfig) error {
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
