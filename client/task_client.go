package client

import (
	"log"
	"strconv"

	"github.com/arbrix/go-test/api"
)

var _ = log.Print

type TaskClient struct {
	Host string
}

func (tc *TaskClient) CreateTask(title, description string, priority int) (api.Task, error) {
	var respTask api.Task
	task := api.Task{Title: title, Description: description, Priority: priority}

	url := tc.Host + "/task"
	r, err := makeRequest("POST", url, task)
	if err != nil {
		return respTask, err
	}
	err = processResponseEntity(r, &respTask, 201)
	return respTask, err
}

func (tc *TaskClient) GetAllTasks() ([]api.Task, error) {
	var respTasks []api.Task

	url := tc.Host + "/task"
	r, err := makeRequest("GET", url, nil)
	if err != nil {
		return respTasks, err
	}
	err = processResponseEntity(r, &respTasks, 200)
	return respTasks, err
}

func (tc *TaskClient) GetTask(id int64) (api.Task, error) {
	var respTask api.Task

	url := tc.Host + "/task/" + strconv.FormatInt(int64(id), 10)
	r, err := makeRequest("GET", url, nil)
	if err != nil {
		return respTask, err
	}
	err = processResponseEntity(r, &respTask, 200)
	return respTask, err
}

func (tc *TaskClient) UpdateTask(task api.Task) (api.Task, error) {
	var respTask api.Task

	url := tc.Host + "/task/" + strconv.FormatInt(int64(task.ID), 10)
	r, err := makeRequest("PUT", url, task)
	if err != nil {
		return respTask, err
	}
	err = processResponseEntity(r, &respTask, 200)
	return respTask, err
}

func (tc *TaskClient) DeleteTask(id int64) error {
	url := tc.Host + "/task/" + strconv.FormatInt(int64(id), 10)
	r, err := makeRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	return processResponse(r, 204)
}
