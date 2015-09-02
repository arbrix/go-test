package task

// CreateTaskForm is used when creating a task.
type CreateTaskForm struct {
	Title       string `json:"title" form:"taskTitle" binding:"required"`
	Description string `json:"desc" form:"taskDesc" binding:"required"`
	Priority    int    `json:"priority" form:"taskPriority" binding:"required"`
}

// UpdateTaskForm is used when updating a Task.
type UpdateTaskForm struct {
	Title       string `json:"title" form:"taskTitle"`
	Description string `json:"desc" form:"taskDesc"`
	Priority    int    `json:"priority" form:"taskPriority"`
	IsCompleted bool   `json:"completed" form:"taskCompleted"`
	IsDeleted   bool   `json:"deleted" form:"taskDeleted"`
}
