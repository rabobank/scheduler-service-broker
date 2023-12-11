package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/rabobank/scheduler-service-broker/model"
	"github.com/rabobank/scheduler-service-broker/util"
	"time"
)

type Schedule struct {
	Guid            string
	SchedulableGuid string
	ExpressionType  string
	Expression      string
	Enabled         int
}

func (schedule Schedule) String() string {
	return fmt.Sprintf("Guid:%s SchedulableGuid:%s, ExpressionType:%s, Expression:%s, Enabled: %d", schedule.Guid, schedule.SchedulableGuid, schedule.ExpressionType, schedule.Expression, schedule.Enabled)
}

func InsertJobSchedule(scheduleCreateRequest model.ScheduleRequest) (string, Job, error) {
	var err error
	var linkedJob Job
	var newGuid string
	db := GetDB()
	defer func() { _ = db.Close() }()
	var jobs []Job
	if jobs, err = GetJobs(scheduleCreateRequest.SpaceGUID, ""); err != nil {
		return newGuid, linkedJob, err
	} else {
		if jobs != nil && len(jobs) == 0 {
			return newGuid, linkedJob, errors.New(fmt.Sprintf("no jobs found for spaceguid %s", scheduleCreateRequest.SpaceGUID))
		}
		jobFound := false
		for _, job := range jobs {
			if job.Name == scheduleCreateRequest.Name {
				jobFound = true
				linkedJob = job
			}
		}
		if !jobFound {
			return newGuid, linkedJob, errors.New(fmt.Sprintf("no jobs found for spaceguid %s and jobname %s", scheduleCreateRequest.SpaceGUID, scheduleCreateRequest.Name))
		}
		newGuid = util.GenerateGUID()
		var result sql.Result
		if result, err = db.Exec("insert into schedules(guid,schedulable_guid,expression,expression_type,enabled) values(?,?,?,?,?)", newGuid, linkedJob.Guid, scheduleCreateRequest.CronExpression, "cron", 1); err != nil {
			return newGuid, linkedJob, errors.New(fmt.Sprintf("failed to insert job schedule, error: %s", err))
		} else {
			if result != nil {
				if inserted, _ := result.RowsAffected(); inserted == 0 {
					return newGuid, linkedJob, errors.New("no job schedule rows inserted")
				}
			}
			fmt.Printf("inserted job schedule with guid %s\n", newGuid)
			return newGuid, linkedJob, err
		}
	}
}

func GetJobSchedules(spaceguid string) ([]model.JobSchedule, error) {
	var err error
	result := make([]model.JobSchedule, 0)
	db := GetDB()
	defer func() { _ = db.Close() }()
	var rows *sql.Rows
	if rows, err = db.Query("select j.appguid,j.name,j.command,s.expression,s.guid from schedules s, jobs j where s.schedulable_guid=j.guid and s.expression!=\"none\" and j.spaceguid=?", spaceguid); err != nil {
		fmt.Printf("failed to query the job schedules, err: %s\n", err)
		return nil, err
	} else {
		result = jobschedules2array(rows)
	}
	return result, nil
}

func jobschedules2array(rows *sql.Rows) []model.JobSchedule {
	result := make([]model.JobSchedule, 0)
	if rows != nil {
		defer func() { _ = rows.Close() }()
		var appguid, jobname, command, expression, scheduleguid string
		for rows.Next() {
			if err := rows.Scan(&appguid, &jobname, &command, &expression, &scheduleguid); err != nil {
				fmt.Printf("failed to scan the schedule row, error:%s\n", err)
			} else {
				appName := "<unknown"
				if app, err := util.CfClient.GetAppByGuid(appguid); err == nil {
					appName = app.Name
				}
				result = append(result, model.JobSchedule{
					AppName:        appName,
					Name:           jobname,
					Command:        command,
					CronExpression: expression,
					ScheduleGuid:   scheduleguid,
				})
			}
		}
	}
	return result
}

func DeleteJobSchedule(jobname, scheduleguid, spaceguid string) error {
	db := GetDB()
	defer func() { _ = db.Close() }()
	if result, err := db.Exec("delete from schedules where guid=? and schedulable_guid in (select guid from jobs where name=? and spaceguid=?)", scheduleguid, jobname, spaceguid); err != nil {
		fmt.Printf("failed to delete schedule, error: %s\n", err)
		return err
	} else {
		numDeletes, _ := result.RowsAffected()
		if numDeletes == 0 {
			return errors.New(fmt.Sprintf("schedule with guid %s and jobname %s does not exist, no rows deleted", scheduleguid, jobname))
		}
		return nil
	}
}

func InsertCallSchedule(scheduleCreateRequest model.ScheduleRequest) (string, Call, error) {
	var err error
	var linkedCall Call
	var newGuid string
	db := GetDB()
	defer func() { _ = db.Close() }()
	var calls []Call
	if calls, err = GetCalls(scheduleCreateRequest.SpaceGUID, ""); err != nil {
		return newGuid, linkedCall, err
	} else {
		if calls != nil && len(calls) == 0 {
			return newGuid, linkedCall, errors.New(fmt.Sprintf("no calls found for spaceguid %s", scheduleCreateRequest.SpaceGUID))
		}
		callFound := false
		for _, call := range calls {
			if call.Name == scheduleCreateRequest.Name {
				callFound = true
				linkedCall = call
			}
		}
		if !callFound {
			return newGuid, linkedCall, errors.New(fmt.Sprintf("no calls found for spaceguid %s and callname %s", scheduleCreateRequest.SpaceGUID, scheduleCreateRequest.Name))
		}
		newGuid = util.GenerateGUID()
		var result sql.Result
		if result, err = db.Exec("insert into schedules(guid,schedulable_guid,expression_type,expression,enabled) values(?,?,?,?,?)", newGuid, linkedCall.Guid, "cron", scheduleCreateRequest.CronExpression, 1); err != nil {
			return newGuid, linkedCall, errors.New(fmt.Sprintf("failed to insert call schedule, error: %s", err))
		} else {
			if result != nil {
				if inserted, _ := result.RowsAffected(); inserted == 0 {
					return newGuid, linkedCall, errors.New("no call schedule rows inserted")
				}
			}
			fmt.Printf("inserted call schedule with guid %s\n", newGuid)
			return newGuid, linkedCall, err
		}
	}
}

func GetCallSchedules(spaceguid string) ([]model.CallSchedule, error) {
	var err error
	result := make([]model.CallSchedule, 0)
	db := GetDB()
	defer func() { _ = db.Close() }()
	var rows *sql.Rows
	if rows, err = db.Query("select c.appguid,c.name,c.url,s.expression,s.guid from schedules s, calls c where s.schedulable_guid=c.guid and s.expression!=\"none\" and c.spaceguid=?", spaceguid); err != nil {
		fmt.Printf("failed to query the call schedules, err: %s\n", err)
		return nil, err
	} else {
		result = callschedules2array(rows)
	}
	return result, nil
}

func callschedules2array(rows *sql.Rows) []model.CallSchedule {
	result := make([]model.CallSchedule, 0)
	if rows != nil {
		defer func() { _ = rows.Close() }()
		var appguid, callname, url, expression, scheduleguid string
		for rows.Next() {
			if err := rows.Scan(&appguid, &callname, &url, &expression, &scheduleguid); err != nil {
				fmt.Printf("failed to scan the schedule row, error:%s\n", err)
			} else {
				appName := "<unknown"
				if app, err := util.CfClient.GetAppByGuid(appguid); err == nil {
					appName = app.Name
				}
				result = append(result, model.CallSchedule{
					AppName:        appName,
					Name:           callname,
					Url:            url,
					CronExpression: expression,
					ScheduleGuid:   scheduleguid,
				})
			}
		}
	}
	return result
}

func DeleteCallSchedule(callname, scheduleguid, spaceguid string) error {
	db := GetDB()
	defer func() { _ = db.Close() }()
	if result, err := db.Exec("delete from schedules where guid=? and schedulable_guid in (select guid from calls where name=? and spaceguid=?)", scheduleguid, callname, spaceguid); err != nil {
		fmt.Printf("failed to delete schedule, error: %s\n", err)
		return err
	} else {
		numDeletes, _ := result.RowsAffected()
		if numDeletes == 0 {
			return errors.New(fmt.Sprintf("schedule with guid %s and callname %s does not exist, no rows deleted", scheduleguid, callname))
		}
		return nil
	}
}

func GetSchedulableJobs() ([]model.SchedulableJob, error) {
	var err error
	result := make([]model.SchedulableJob, 0)
	db := GetDB()
	defer func() { _ = db.Close() }()
	var rows *sql.Rows
	if rows, err = db.Query("select s.guid,j.appguid,j.spaceguid,s.expression,j.name as jobname,j.command,j.memoryinmb,j.diskinmb from schedules s, jobs j where s.expression_type=\"cron\" and s.expression!=\"none\" and s.schedulable_guid=j.guid and s.enabled=1"); err != nil {
		fmt.Printf("failed to query the job schedules/jobs, err: %s\n", err)
		return nil, err
	} else {
		if rows != nil {
			defer func() { _ = rows.Close() }()
			var scheduleguid, expression, jobname, appguid, spaceguid, command string
			var memoryInMB, diskInMB int
			for rows.Next() {
				if err = rows.Scan(&scheduleguid, &appguid, &spaceguid, &expression, &jobname, &command, &memoryInMB, &diskInMB); err != nil {
					fmt.Printf("failed to scan the schedules/jobs row:%s\n", err)
				} else {
					result = append(result, model.SchedulableJob{
						ScheduleGuid: scheduleguid,
						Expression:   expression,
						JobName:      jobname,
						AppGuid:      appguid,
						SpaceGuid:    spaceguid,
						Command:      command,
						MemoryInMB:   memoryInMB,
						DiskInMB:     diskInMB,
					})
				}
			}
		}
	}
	return result, nil
}

func GetSchedulableCalls() ([]model.SchedulableCall, error) {
	var err error
	result := make([]model.SchedulableCall, 0)
	db := GetDB()
	defer func() { _ = db.Close() }()
	var rows *sql.Rows
	if rows, err = db.Query("select s.guid,c.appguid,c.spaceguid, s.expression,c.name as callname,c.url, c.authheader from schedules s, calls c where s.expression_type=\"cron\" and s.expression!=\"none\" and s.schedulable_guid=c.guid and s.enabled=1"); err != nil {
		fmt.Printf("failed to query the call schedules/calls, err: %s\n", err)
		return nil, err
	} else {
		if rows != nil {
			defer func() { _ = rows.Close() }()
			var scheduleguid, expression, callname, appguid, spaceguid, url, authheader string
			for rows.Next() {
				if err = rows.Scan(&scheduleguid, &appguid, &spaceguid, &expression, &callname, &url, &authheader); err != nil {
					fmt.Printf("failed to scan the schedules/calls row:%s\n", err)
				} else {
					result = append(result, model.SchedulableCall{
						ScheduleGuid: scheduleguid,
						Expression:   expression,
						CallName:     callname,
						AppGuid:      appguid,
						SpaceGuid:    spaceguid,
						Url:          url,
						AuthHeader:   authheader,
					})
				}
			}
		}
	}
	return result, nil
}

// DeleteScheduleByAge - When you use the run-job or run-call commands, we have to create a schedule of type "execute", these have to be cleaned up every now and then
func DeleteScheduleByAge(maxDays int64) error {
	var err error
	nanoseconds := 1000000000 * 3600 * 24 * maxDays
	cutOffDate := time.Now().Add(time.Duration(-nanoseconds))
	fmt.Printf("cleaning schedules table for rows older than %d days (older than %s)\n", maxDays, cutOffDate.Format(time.RFC3339))
	db := GetDB()
	defer func() { _ = db.Close() }()
	var result sql.Result
	if result, err = db.Exec("delete from schedules where expression_type=\"execute\" and guid in (select schedule_guid from histories where execution_start_time<?)", cutOffDate); err != nil {
		fmt.Printf("failed to clean up schedules: %s\n", err)
		return err
	}
	if numDeleted, _ := result.RowsAffected(); numDeleted > 0 {
		fmt.Printf("deleted %d rows from schedules table\n", numDeleted)
	}
	return nil
}
