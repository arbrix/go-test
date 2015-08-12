package web

import (
	_ "database/sql"
	"fmt"

	"github.com/arbrix/go-test/controllers"
	"github.com/arbrix/go-test/web/middleware"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/tommy351/gin-cors"
)

type Config struct {
	ListenAddress string
	DatabaseUri   string
}

func (cfg Config) String() string {
	return cfg.ListenAddress + "; " + cfg.DatabaseUri
}

type TaskService struct {
}

func (s *TaskService) getDb(cfg Config) (gorm.DB, error) {
	fmt.Println(cfg.DatabaseUri)
	return gorm.Open("mysql", cfg.DatabaseUri)
}

func (s *TaskService) Run(cfg Config) error {
	db, err := s.getDb(cfg)
	if err != nil {
		return err
	}
	db.SingularTable(true)
	db.LogMode(true)

	taskController := &controllers.TaskController{Db: db}

	r := gin.Default()
	r.Use(cors.Middleware(cors.Options{}))
	r.Use(middleware.CheckHeader())

	r.GET("/task", taskController.GetAllTasks)
	r.GET("/task/:id", taskController.GetTask)
	r.POST("/task", taskController.CreateTask)
	r.PUT("/task/:id", taskController.UpdateTask)
	r.DELETE("/task/:id", taskController.DeleteTask)
	r.DELETE("/task-del/:id", taskController.RealDeleteTask)

	r.Run(cfg.ListenAddress)

	return nil
}
