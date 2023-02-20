package model

type JobListResponse struct {
	Jobs []Job
}

type Job struct {
	JobName string `json:"jobname"`
	AppName string `json:"appname"`
	Command string `json:"command"`
}
