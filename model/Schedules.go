package model

type ScheduleRequest struct {
	SpaceGUID      string `json:"spaceguid"`
	Name           string `json:"name"`
	CronExpression string `json:"cronexpression"`
	ExpressionType string
}

type JobScheduleListResponse struct {
	JobSchedules []JobSchedule
}

type JobSchedule struct {
	AppName        string `json:"appname"`
	Name           string `json:"name"`
	Command        string `json:"command"`
	CronExpression string `json:"cronexpression"`
	ScheduleGuid   string `json:"scheduleguid"`
}

type CallScheduleListResponse struct {
	CallSchedules []CallSchedule
}

type CallSchedule struct {
	AppName        string `json:"appname"`
	Name           string `json:"name"`
	Url            string `json:"Url"`
	CronExpression string `json:"cronexpression"`
	ScheduleGuid   string `json:"scheduleguid"`
}

type SchedulableJob struct {
	ScheduleGuid string
	Expression   string
	JobName      string
	AppGuid      string
	SpaceGuid    string
	Command      string
}

type SchedulableCall struct {
	ScheduleGuid string
	Expression   string
	CallName     string
	AppGuid      string
	SpaceGuid    string
	Url          string
	AuthHeader   string
}
