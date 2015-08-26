package v1

import (
	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/gin"

	"github.com/arbrix/go-test/api/response"
	"github.com/arbrix/go-test/config"
	"github.com/arbrix/go-test/service/taskService"
)

// @Title Tasks
// @Description Tasks's router group.
func Tasks(parentRoute *gin.RouterGroup) {
	route := parentRoute.Group("/tasks")
	route.Use(jwt.Auth(config.SecretKey))
	route.POST("", createTask)
	route.GET("/:id", retrieveTask)
	route.GET("", retrieveTasks)
	route.PUT("/:id", updateTask)
	route.DELETE("/:id", deleteTask)
	parentRoute.DELETE("/task-del/:id", realDeleteTask)
}

// @Title createTask
// @Description Create a task.
// @Accept  json
// @Param   taskTitle       form   string  true        "Task Title."
// @Param   taskDesc        form   string  true        "Task Description."
// @Param   taskPriorit     form   int     true        "Task Priority."
// @Success 201 {object} response.BasicResponse "Task is created successfully"
// @Failure 401 {object} response.BasicResponse "Authentication required"
// @Failure 404 {object} response.BasicResponse "User not logged in."
// @Failure 500 {object} response.BasicResponse "Task is not created."
// @Resource /tasks
// @Router /tasks [post]
func createTask(c *gin.Context) {
	status, err := taskService.CreateTask(c)
	messageTypes := &response.MessageTypes{
		OK:                  "creation.done",
		Unauthorized:        "login.error.fail",
		NotFound:            "creation.error.fail",
		InternalServerError: "creation.error.fail",
	}
	messages := &response.Messages{OK: "Task is created successfully."}
	response.JSON(c, status, messageTypes, messages, err)
}

// @Title retrieveTask
// @Description Retrieve a task.
// @Accept  json
// @Param   id        path    int     true        "Task ID"
// @Success 200 {object} model.Task "OK"
// @Failure 404 {object} response.BasicResponse "Not found"
// @Resource /tasks
// @Router /tasks/{id} [get]
func retrieveTask(c *gin.Context) {
	task, status, err := taskService.RetrieveTask(c)
	if err == nil {
		c.JSON(status, gin.H{"task": task})
	} else {
		messageTypes := &response.MessageTypes{
			NotFound:     "task.error.notFound",
			Unauthorized: "task.error.unauthorized",
		}
		response.ErrorJSON(c, status, messageTypes, err)
	}

}

// @Title retrieveTasks
// @Description Retrieve task array.
// @Accept  json
// @Success 200 {array} model.Task "OK"
// @Resource /tasks
// @Router /tasks [get]
func retrieveTasks(c *gin.Context) {
	tasks := taskService.RetrieveTasks(c)
	c.JSON(200, gin.H{"tasks": tasks})
}

// @Title updateTask
// @Description Update a task.
// @Accept  json
// @Param   id        path    int     true        "Task ID"
// @Param   taskTitle       form   string  true        "Task Title."
// @Param   taskDesc        form   string  true        "Task Description."
// @Param   taskPriorit     form   int     true        "Task Priority."
// @Param   taskCompleted   form   bool  true        "Task Complete mark."
// @Param   taskDeleted     form   bool  true        "Task Delete mark."
// @Success 200 {object} model.Task "OK"
// @Failure 401 {object} response.BasicResponse "Authentication required"
// @Failure 404 {object} response.BasicResponse "Not found"
// @Failure 500 {object} response.BasicResponse "Task is not updated."
// @Resource /tasks
// @Router /tasks/{id} [put]
func updateTask(c *gin.Context) {
	task, status, err := taskService.UpdateTask(c)
	if err == nil {
		c.JSON(status, gin.H{"task": task})
	} else {
		messageTypes := &response.MessageTypes{
			Unauthorized:        "task.error.unauthorized",
			NotFound:            "task.error.notFound",
			InternalServerError: "task.error.internalServerError",
		}
		response.ErrorJSON(c, status, messageTypes, err)
	}
}

// @Title deleteTask
// @Description Delete a task.
// @Accept  json
// @Param   id        path    int     true        "Task ID"
// @Success 200 {object} response.BasicResponse
// @Failure 401 {object} response.BasicResponse "Authentication required"
// @Failure 404 {object} response.BasicResponse "Not found"
// @Failure 500 {object} response.BasicResponse "Task is not deleted."
// @Resource /tasks
// @Router /tasks/{id} [delete]
func deleteTask(c *gin.Context) {
	_, status, err := taskService.MarkAsDeleted(c)
	if err == nil {
		c.JSON(status, response.BasicResponse{})
	} else {
		messageTypes := &response.MessageTypes{
			Unauthorized:        "user.error.unauthorized",
			NotFound:            "task.error.notFound",
			InternalServerError: "setting.leaveOurService.fail",
		}
		response.ErrorJSON(c, status, messageTypes, err)
	}
}

// @Title realDeleteTask
// @Description Delete a task from DB.
// @Accept  json
// @Param   id        path    int     true        "Task ID"
// @Success 200 {object} response.BasicResponse
// @Failure 401 {object} response.BasicResponse "Authentication required"
// @Failure 404 {object} response.BasicResponse "Not found"
// @Failure 500 {object} response.BasicResponse "Task is not deleted."
// @Resource /tasks
// @Router /tasks/del/{id} [delete]
func realDeleteTask(c *gin.Context) {
	status, err := taskService.DeleteTask(c)
	if err == nil {
		c.JSON(status, response.BasicResponse{})
	} else {
		messageTypes := &response.MessageTypes{
			Unauthorized:        "user.error.unauthorized",
			NotFound:            "task.error.notFound",
			InternalServerError: "setting.leaveOurService.fail",
		}
		response.ErrorJSON(c, status, messageTypes, err)
	}
}
