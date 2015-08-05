package service

import (
	"github.com/arbrix/go-test/api"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Config struct {
	SvcHost    string
	DbUser     string
	DbPassword string
	DbHost     string
	DbName     string
}

type TaskService struct {
}

func (s *TaskService) getDb(cfg Config) (gorm.DB, error) {
	connectionString := cfg.DbUser + ":" + cfg.DbPassword + "@tcp(" + cfg.DbHost + ":3306)/" + cfg.DbName + "?charset=utf8&parseTime=True"

	return gorm.Open("mysql", connectionString)
}

func (s *TaskService) Migrate(cfg Config) error {
	db, err := s.getDb(cfg)
	if err != nil {
		return err
	}
	db.SingularTable(true)

	db.AutoMigrate(&api.Task{})
	return nil
}
func (s *TaskService) Run(cfg Config) error {
	db, err := s.getDb(cfg)
	if err != nil {
		return err
	}
	db.SingularTable(true)

	taskResource := &TaskResource{db: db}

	r := gin.Default()

	auth := r.Group("/")
	auth.Use(CheckHeader())
	{
		r.GET("/task", taskResource.GetAllTasks)
		r.GET("/task/:id", taskResource.GetTask)
		r.POST("/task", taskResource.CreateTask)
		r.PUT("/task/:id", taskResource.UpdateTask)
		r.PATCH("/task/:id", taskResource.PatchTask)
		r.DELETE("/task/:id", taskResource.DeleteTask)
	}

	r.Run(cfg.SvcHost)

	return nil
}
