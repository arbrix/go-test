package controllers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/arbrix/go-test/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type TaskController struct {
	db gorm.DB
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	var task models.Task

	if c.Bind(&task) == nil {
		fmt.Println(task.String())
		c.JSON(400, models.NewError("problem decoding body"))
		return
	}
	task.CreatedAt = time.Now()
	fmt.Println(task.String())

	tc.db.Save(&task)

	c.JSON(201, task)
}

func (tc *TaskController) GetAllTasks(c *gin.Context) {
	var tasks []models.Task

	tc.db.Order("CreatedAt desc").Find(&tasks)

	c.JSON(200, tasks)
}

func (tc *TaskController) GetTask(c *gin.Context) {
	id, err := tc.getId(c)
	if err != nil {
		c.JSON(400, models.NewError("problem decoding id sent"))
		return
	}

	var task models.Task

	if tc.db.First(&task, id).RecordNotFound() {
		c.JSON(404, gin.H{"error": "not found"})
	} else {
		c.JSON(200, task)
	}
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	id, err := tc.getId(c)
	if err != nil {
		c.JSON(400, models.NewError("problem decoding id sent"))
		return
	}

	var task models.Task

	if c.Bind(&task) == nil {
		fmt.Println(task.String())
		c.JSON(400, models.NewError("problem decoding body"))
		return
	}
	task.ID = int64(id)
	task.UpdatedAt = time.Now()
	if task.IsCompleted == true {
		task.CompletedAt = time.Now()
	}
	fmt.Println(task.String())
	var existing models.Task

	if tc.db.First(&existing, id).RecordNotFound() {
		c.JSON(404, models.NewError("not found"))
	} else {
		tc.db.Save(&task)
		c.JSON(200, task)
	}

}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	id, err := tc.getId(c)
	if err != nil {
		c.JSON(400, models.NewError("problem decoding id sent"))
		return
	}

	var task models.Task
	task.IsDeleted = true

	if tc.db.First(&task, id).RecordNotFound() {
		c.JSON(404, models.NewError("not found"))
	} else {
		tc.db.Save(&task)
		c.Data(204, "application/json", make([]byte, 0))
	}
}

func (tc *TaskController) getId(c *gin.Context) (int64, error) {
	idStc := c.Params.ByName("id")
	id, err := strconv.Atoi(idStc)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	return int64(id), nil
}
