package domain

import (
	"time"
)

type TaskStatus string

const (
	_                         = iota
	TaskIncomplete TaskStatus = "Incomplete"
	TaskCompleted  TaskStatus = "Completed"
)

type Model struct {
	UID       string    `json:"uid"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Task struct {
	Model
	Content  string     `json:"content"`
	Status   TaskStatus `json:"status"`
	Priority int        `json:"priority"`
}
