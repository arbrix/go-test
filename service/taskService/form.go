package taskService

// CreateTaskForm is used when creating a task.
type CreateTaskForm struct {
	Title       string `form:"taskTitle" binding:"required"`
	Description string `form:"taskDesc" binding:"required"`
	Priority    int    `form:"taskPriority" binding:"required"`
}

// UpdateTaskForm is used when updating a Task.
type UpdateTaskForm struct {
	Title       string `form:"taskTitle"`
	Description string `form:"taskDesc"`
	Priority    int    `form:"taskPriority"`
	IsCompleted bool   `form:"taskCompleted"`
	IsDeleted   bool   `form:"taskDeleted"`
}
