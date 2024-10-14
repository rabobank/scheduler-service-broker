package controllers

import (
	"fmt"
	"github.com/rabobank/scheduler-service-broker/conf"
	"github.com/rabobank/scheduler-service-broker/cron"
	"github.com/rabobank/scheduler-service-broker/db"
	"github.com/rabobank/scheduler-service-broker/model"
	"github.com/rabobank/scheduler-service-broker/util"
	"net/http"
	"time"
)

func JobCreate(w http.ResponseWriter, r *http.Request) {
	var err error
	if isValid, userId, req := ValidateRequest(w, r); isValid {
		if req.Name == "" {
			util.WriteHttpResponse(w, http.StatusBadRequest, "the name of the job was empty")
		} else {
			if _, err = conf.CfClient.Applications.Get(conf.CfCtx, req.AppGUID); err != nil {
				util.WriteHttpResponse(w, http.StatusBadRequest, fmt.Sprintf("app with guid %s not found: %s", req.AppGUID, err))
			} else {
				if !util.IsAppBoundToSchedulerService(req.AppGUID) {
					util.WriteHttpResponse(w, http.StatusBadRequest, fmt.Sprintf("app with guid %s is not bound to an instance of scheduler", req.AppGUID))
				} else {
					// try to insert the job
					if jobguid, err := db.InsertJob(db.Job{AppGuid: req.AppGUID, SpaceGuid: req.SpaceGUID, Name: req.Name, Command: req.Command, MemoryInMB: req.MemoryInMB, DiskInMB: req.DiskInMB}); err != nil {
						util.WriteHttpResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to create job: %s", err))
					} else {
						fmt.Printf("userId %s created job with guid %s for space guid %s\n", userId, jobguid, req.SpaceGUID)
						util.WriteHttpResponse(w, http.StatusCreated, fmt.Sprintf("job created"))
					}
				}
			}
		}
	}
}

func JobRun(w http.ResponseWriter, r *http.Request) {
	if isValid, userId, req := ValidateRequest(w, r); isValid {
		if req.Name == "" {
			util.WriteHttpResponse(w, http.StatusBadRequest, "the name of the job was empty")
		} else {
			if existingJobs, err := db.GetJobs(req.SpaceGUID, req.Name); err == nil && len(existingJobs) > 0 {
				// there should only be a result of length 1:
				job := existingJobs[0]
				go func() {
					cron.DoJob(time.Now(), model.SchedulableJob{JobName: job.Name, AppGuid: job.AppGuid, SpaceGuid: job.SpaceGuid, Command: job.Command, MemoryInMB: job.MemoryInMB, DiskInMB: job.DiskInMB})
				}()
				fmt.Printf("userId %s ran job for space guid %s\n", userId, req.SpaceGUID)
				util.WriteHttpResponse(w, http.StatusOK, fmt.Sprintf("run scheduled for job %s", req.Name))
			} else {
				util.WriteHttpResponse(w, http.StatusBadRequest, fmt.Sprintf("job %s was not found", req.Name))
			}
		}
	}
}

func JobGet(w http.ResponseWriter, r *http.Request) {
	if isValid, _, req := ValidateRequest(w, r); isValid {
		if result, err := db.GetJobs(req.SpaceGUID, ""); err != nil {
			util.WriteHttpResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to list jobs: %s", err))
		} else {
			var jobs = make([]model.Job, 0)
			for _, job := range result {
				appName := "<unknown>"
				if app, err := conf.CfClient.Applications.Get(conf.CfCtx, job.AppGuid); err == nil {
					appName = app.Name
				}
				jobs = append(jobs, model.Job{JobName: job.Name, AppName: appName, Command: job.Command, MemoryInMB: job.MemoryInMB, DiskInMB: job.DiskInMB})
			}
			util.WriteHttpResponse(w, http.StatusOK, model.JobListResponse{Jobs: jobs})
		}
	}
}

func JobDelete(w http.ResponseWriter, r *http.Request) {
	var err error
	if isValid, userId, req := ValidateRequest(w, r); isValid {
		if req.Name == "" {
			util.WriteHttpResponse(w, http.StatusBadRequest, "the name of the job was empty")
		} else {
			fmt.Printf("deleting job with jobname %s for userId %s for spaceguid %s from client %s...\n", req.Name, userId, req.SpaceGUID, r.RemoteAddr)
			if err = db.DeleteJobBySpaceGuidAndJobname(req.SpaceGUID, req.Name); err != nil {
				util.WriteHttpResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete jobs with name %s in spaceguid %s: %s", req.Name, req.SpaceGUID, err))
			} else {
				cron.DeleteJobByNameAndSpaceGuid(req.Name, req.SpaceGUID)
				util.WriteHttpResponse(w, http.StatusOK, fmt.Sprintf("job %s deleted", req.Name))
			}
		}
	}
}
