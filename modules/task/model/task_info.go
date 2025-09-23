package model

import (
	"time"
	"github.com/vera-byte/vgo/v"
)

const TableNameTaskInfo = "task_info"

// TaskInfo mapped from table <task_info>
type TaskInfo struct {
	*v.Model
	JobId       string     `json:"jobId"`
	RepeatConf  string     `json:"repeatConf"`
	Name        string     `json:"name"`
	Cron        string     `json:"cron"`
	Limit       int        `json:"limit"`
	Every       int        `json:"every"`
	Remark      string     `json:"remark"`
	Status      int        `json:"status"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Data        string     `json:"data"`
	Service     string     `json:"service"`
	Type        int        `json:"type"`
	NextRunTime time.Time `json:"nextRunTime"`
	TaskType    int        `json:"taskType"`
}

// TableName TaskInfo's table name
func (*TaskInfo) TableName() string {
	return TableNameTaskInfo
}

// GroupName TaskInfo's table group
func (*TaskInfo) GroupName() string {
	return "default"
}

// NewTaskInfo create a new TaskInfo
func NewTaskInfo() *TaskInfo {
	return &TaskInfo{
		Model: v.NewModel(),
	}
}
