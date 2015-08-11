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
	_ = client.DeleteTask(task.ID)
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
	_ = client.DeleteTask(task.ID)
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

	_, err = client.GetTask(id)
	if err == nil {
		t.Error(err)
	}
}
