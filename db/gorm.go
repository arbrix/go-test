package db

import (
	"fmt"
	"github.com/arbrix/go-test/config"
	"github.com/arbrix/go-test/model"
	"github.com/arbrix/go-test/util/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var ORM gorm.DB

// GormInit init gorm ORM.
func GormInit(cfg config.Config) gorm.DB {
	fmt.Println(cfg.DatabaseUri)
	db, err := gorm.Open("mysql", cfg.DatabaseUri)
	db.DB()

	// Then you could invoke `*sql.DB`'s functions with it
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// Disable table name's pluralization
	db.SingularTable(true)
	if config.Environment == "DEVELOPMENT" {
		db.LogMode(true)
		db.AutoMigrate(&model.User{}, &model.Connection{}, &model.Task{})
	}
	log.CheckError(err)

	return db
}
