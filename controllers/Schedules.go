package controllers

import (
	"fmt"
	"github.com/rabobank/scheduler-service-broker/cron"
	"github.com/rabobank/scheduler-service-broker/db"
	"github.com/rabobank/scheduler-service-broker/model"
	"github.com/rabobank/scheduler-service-broker/util"
	"net/http"
)

func JobScheduleCreate(w http.ResponseWriter, r *http.Request) {
	if isValid, userId, req := ValidateRequest(w, r); isValid {
		// try to insert the job-schedule
		//req.ExpressionType = "cron_expression"
		if scheduleGuid, job, err := db.InsertJobSchedule(model.ScheduleRequest{SpaceGUID: req.SpaceGUID, Name: req.Name, CronExpression: req.CronExpression, ExpressionType: req.ExpressionType}); err != nil {
			util.WriteHttpResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to create schedule: %s", err))
		} else {
			fmt.Printf("userId %s created schedule for job %s for space guid %s\n", userId, job.Name, req.SpaceGUID)
			cron.AddJob(model.SchedulableJob{ScheduleGuid: scheduleGuid, Expression: req.CronExpression, JobName: req.Name, AppGuid: job.AppGuid, SpaceGuid: req.SpaceGUID, Command: job.Command})
			util.WriteHttpResponse(w, http.StatusCreated, fmt.Sprintf("schedule created"))
		}
	}
}

func JobScheduleGet(w http.ResponseWriter, r *http.Request) {
	if isValid, _, req := ValidateRequest(w, r); isValid {
		if result, err := db.GetJobSchedules(req.SpaceGUID); err != nil {
			util.WriteHttpResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to list job schedules: %s", err))
		} else {
			util.WriteHttpResponse(w, http.StatusOK, model.JobScheduleListResponse{JobSchedules: result})
		}
	}
}

func JobScheduleDelete(w http.ResponseWriter, r *http.Request) {
	var err error
	if isValid, userId, req := ValidateRequest(w, r); isValid {
		fmt.Printf("deleting schedule %s for userId %s for spaceguid %s from client %s...\n", req.Name, userId, req.SpaceGUID, r.RemoteAddr)
		if err = db.DeleteJobSchedule(req.Name, req.ScheduleGuid, req.SpaceGUID); err != nil {
			util.WriteHttpResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete schedules for job name %s in spaceguid %s: %s", req.Name, req.SpaceGUID, err))
		} else {
			cron.DeleteJobByScheduleGuid(req.ScheduleGuid)
			util.WriteHttpResponse(w, http.StatusOK, fmt.Sprintf("schedule for job name %s deleted", req.Name))
		}
	}
}

func CallScheduleCreate(w http.ResponseWriter, r *http.Request) {
	if isValid, userId, req := ValidateRequest(w, r); isValid {
		// try to insert the call-schedule
		if scheduleGuid, call, err := db.InsertCallSchedule(model.ScheduleRequest{SpaceGUID: req.SpaceGUID, Name: req.Name, CronExpression: req.CronExpression, ExpressionType: req.ExpressionType}); err != nil {
			util.WriteHttpResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to create schedule: %s", err))
		} else {
			fmt.Printf("userId %s created schedule for call %s for space guid %s\n", userId, call.Name, req.SpaceGUID)
			cron.AddCall(model.SchedulableCall{ScheduleGuid: scheduleGuid, Expression: req.CronExpression, CallName: req.Name, AppGuid: call.AppGuid, SpaceGuid: req.SpaceGUID, Url: call.Url})
			util.WriteHttpResponse(w, http.StatusCreated, fmt.Sprintf("schedule created"))
		}
	}
}

func CallScheduleGet(w http.ResponseWriter, r *http.Request) {
	if isValid, _, req := ValidateRequest(w, r); isValid {
		if result, err := db.GetCallSchedules(req.SpaceGUID); err != nil {
			util.WriteHttpResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to list call schedules: %s", err))
		} else {
			util.WriteHttpResponse(w, http.StatusOK, model.CallScheduleListResponse{CallSchedules: result})
		}
	}
}

func CallScheduleDelete(w http.ResponseWriter, r *http.Request) {
	var err error
	if isValid, userId, req := ValidateRequest(w, r); isValid {
		fmt.Printf("deleting schedule %s for userId %s for spaceguid %s from client %s...\n", req.Name, userId, req.SpaceGUID, r.RemoteAddr)
		if err = db.DeleteCallSchedule(req.Name, req.ScheduleGuid, req.SpaceGUID); err != nil {
			util.WriteHttpResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete schedules for call name %s in spaceguid %s: %s", req.Name, req.SpaceGUID, err))
		} else {
			cron.DeleteCallByScheduleGuid(req.ScheduleGuid)
			util.WriteHttpResponse(w, http.StatusOK, fmt.Sprintf("schedule for call name %s deleted", req.Name))
		}
	}
}
