package model

type JobListResponse struct {
	Jobs []Job
}

type Job struct {
	JobName    string `json:"jobname"`
	AppName    string `json:"appname"`
	Command    string `json:"command"`
	MemoryInMB int    `json:"memory_in_mb,omitempty" `
	DiskInMB   int    `json:"disk_in_mb,omitempty" `
}
