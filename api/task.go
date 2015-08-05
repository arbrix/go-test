package api

import "time"

type Task struct {
	Id          int64      `json:"id"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Priority    int        `json:"priority"`
	CreatedAt   *time.Time `json:"created"`
	UpdatedAt   *time.Time `json:"updated"`
	CompletedAt *time.Time `json:"completed"`
	IsDeleted   bool       `json:"isDeleted"`
	IsCompleted bool       `json:"isCompeted"`
}
