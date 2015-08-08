package service

import (
	"github.com/arbrix/go-test/api"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/tommy351/gin-cors"
)

type Config struct {
	ListenAddress string
	DatabaseUri   string
}

type TaskService struct {
}

func (s *TaskService) getDb(cfg Config) (gorm.DB, error) {
	return gorm.Open("mysql", cfg.DatabaseUri)
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

func index(c *gin.Context) {
	content := gin.H{"Hello": "World"}
	c.JSON(200, content)
}

func (s *TaskService) Run(cfg Config) error {
	db, err := s.getDb(cfg)
	if err != nil {
		return err
	}
	db.SingularTable(true)

	taskResource := &TaskResource{db: db}

	r := gin.Default()
	r.Use(cors.Middleware(cors.Options{}))

	auth := r.Group("/")
	auth.Use(CheckHeader())
	{
		r.GET("/task", taskResource.GetAllTasks)
		r.GET("/task/:id", taskResource.GetTask)
		r.POST("/task", taskResource.CreateTask)
		r.PUT("/task/:id", taskResource.UpdateTask)
		r.DELETE("/task/:id", taskResource.DeleteTask)
	}
	r.GET("/test", index)

	r.Run(cfg.ListenAddress)

	return nil
}
