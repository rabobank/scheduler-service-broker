package model

import (
	"fmt"
	"time"
)

type HistoryListResponse struct {
	Histories []History
}

type History struct {
	Guid               string
	ScheduledTime      time.Time
	ExecutionStartTime time.Time
	ExecutionEndTime   time.Time
	Message            string
	State              string
	ScheduleGuid       string
	TaskGuid           string
	CreatedAt          time.Time
}

func (h History) String() string {
	return fmt.Sprintf("Guid:%s ScheduledTime:%s, ExecutionStartTime:%s, ExecutionEndTime:%s, Message: %s, ScheduleGuid:%s, TaskGuid:%s", h.Guid, h.ScheduledTime, h.ExecutionStartTime, h.ExecutionEndTime, h.Message, h.ScheduleGuid, h.TaskGuid)
}
