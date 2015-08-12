package tests

import (
	"testing"

	"github.com/arbrix/go-test/utils"
)

func TestCreateTask(t *testing.T) {

	// given
	client := utils.TaskClient{Host: "http://localhost:8080"}

	// when
	task, err := client.CreateTask("foo", "bar", 1)

	//then
	if err != nil {
		t.Error(err)
	}

	if task.Title != "foo" && task.Description != "bar" && task.Priority != 1 {
		t.Error("returned task not right")
	}

	// cleanup
	_ = client.RealDeleteTask(task.ID)
}

func TestGetTask(t *testing.T) {

	// given
	client := utils.TaskClient{Host: "http://localhost:8080"}
	task, _ := client.CreateTask("foo", "bar", 1)
	id := task.ID

	// when
	task, err := client.GetTask(id)

	// then
	if err != nil {
		t.Error(err)
	}

	if task.Title != "foo" && task.Description != "bar" && task.Priority != 1 {
		t.Error("returned task not right")
	}

	// cleanup
	_ = client.RealDeleteTask(task.ID)
}

func TestGetAllTask(t *testing.T) {
	// given
	client := utils.TaskClient{Host: "http://localhost:8080"}
	task1, _ := client.CreateTask("task1", "Desk 1", 1)
	task2, _ := client.CreateTask("task2", "Desk 2", 2)

	// when
	tasks, err := client.GetAllTasks()

	// then
	if err != nil {
		t.Error(err)
	}

	if len(tasks) < 2 {
		t.Error("getting all tasks works not correct")
	}

	// cleanup
	_ = client.RealDeleteTask(task1.ID)
	_ = client.RealDeleteTask(task2.ID)
}

func TestCompliteTask(t *testing.T) {
	// given
	client := utils.TaskClient{Host: "http://localhost:8080"}
	task, _ := client.CreateTask("foo4complite", "bar for complite", 3)
	task.IsCompleted = true
	// when
	task, err := client.UpdateTask(task)
	// then
	if err != nil {
		t.Error(err)
	}
	if task.Title != "foo4complite" && task.Description != "bar for complite" && task.Priority != 3 && task.IsCompleted != true {
		t.Error("returned task not right")
	}
	// cleanup
	_ = client.RealDeleteTask(task.ID)
}

func TestDeleteTask(t *testing.T) {

	// given
	client := utils.TaskClient{Host: "http://localhost:8080"}
	task, _ := client.CreateTask("foo", "bar", 1)
	id := task.ID

	// when
	err := client.DeleteTask(id)

	// then
	if err != nil {
		t.Error(err)
	}
	_ = client.RealDeleteTask(id)
}
