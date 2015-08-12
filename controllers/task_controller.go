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
	Db gorm.DB
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	var task models.Task

	if c.Bind(&task) != nil {
		fmt.Println(task.String())
		c.JSON(400, models.NewError("problem decoding body"))
		return
	}
	task.CreatedAt = time.Now()
	fmt.Println(task.String())

	tc.Db.Save(&task)

	c.JSON(201, task)
}

func (tc *TaskController) GetAllTasks(c *gin.Context) {
	var tasks []models.Task

	tc.Db.Where("isDeleted = ?", 0).Order("created desc").Find(&tasks)

	c.JSON(200, tasks)
}

func (tc *TaskController) GetTask(c *gin.Context) {
	id, err := tc.getId(c)
	if err != nil {
		c.JSON(400, models.NewError("problem decoding id sent"))
		return
	}

	var task models.Task

	if tc.Db.Where("isDeleted = ?", 0).First(&task, id).RecordNotFound() {
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
	tc.Db.First(&task, id)

	if c.Bind(&task) != nil {
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
	existing.IsDeleted = false

	if tc.Db.First(&existing, id).RecordNotFound() {
		c.JSON(404, models.NewError("not found"))
	} else {
		tc.Db.Save(&task)
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

	if tc.Db.Where("isDeleted = ?", 0).First(&task, id).RecordNotFound() {
		c.JSON(404, models.NewError("not found"))
	} else {
		task.IsDeleted = true
		task.UpdatedAt = time.Now()
		tc.Db.Save(&task)
		c.Data(204, "application/json", make([]byte, 0))
	}
}

func (tc *TaskController) RealDeleteTask(c *gin.Context) {
	id, err := tc.getId(c)
	if err != nil {
		c.JSON(400, models.NewError("problem decoding id sent"))
		return
	}

	var task models.Task

	if tc.Db.Where("isDeleted = ?", 0).First(&task, id).RecordNotFound() {
		c.JSON(404, models.NewError("not found"))
	} else {
		tc.Db.Delete(&task)
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
