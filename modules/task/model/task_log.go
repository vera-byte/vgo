package model

import (
	"github.com/vera-byte/vgo/v"
)

const TableNameTaskLog = "task_log"

// TaskLog mapped from table <task_log>
type TaskLog struct {
	*v.Model
	TaskId uint64 `json:"taskId"`
	Status uint8  `json:"status"`
	Detail string `json:"detail"`
}

// TableName TaskLog's table name
func (*TaskLog) TableName() string {
	return TableNameTaskLog
}

// GroupName TaskLog's table group
func (*TaskLog) GroupName() string {
	return "default"
}

// NewTaskLog create a new TaskLog
func NewTaskLog() *TaskLog {
	return &TaskLog{
		Model: v.NewModel(),
	}
}
