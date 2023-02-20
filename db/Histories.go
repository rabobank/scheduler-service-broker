package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/rabobank/scheduler-service-broker/model"
	"github.com/rabobank/scheduler-service-broker/util"
	"time"
)

func InsertHistory(history model.History) (string, error) {
	var err error
	db := GetDB()
	defer db.Close()
	history.Guid = util.GenerateGUID()
	util.PrintfIfDebug("inserting history: %s\n", history)
	var scheduleGuid sql.NullString
	if len(history.ScheduleGuid) != 0 {
		scheduleGuid = sql.NullString{String: history.ScheduleGuid, Valid: true}
	}
	if _, err = db.Exec("insert into histories(guid,scheduled_time,execution_start_time,execution_end_time,state,message,schedule_guid,task_guid) values(?,?,?,?,?,?,?,?)", history.Guid, history.ScheduledTime, history.ExecutionStartTime, history.ExecutionEndTime, history.State, history.Message, scheduleGuid, history.TaskGuid); err != nil {
		fmt.Printf("failed to insert history: %s\n", err)
	} else {
		util.PrintfIfDebug("inserted history with guid %s\n", history.Guid)
	}
	return history.Guid, err
}

func GetJobHistories(spaceguid, name string) ([]model.History, error) {
	var err error
	if spaceguid == "" {
		spaceguid = "%"
	}
	if name == "" {
		name = "%"
	}
	result := make([]model.History, 0)
	db := GetDB()
	defer db.Close()
	var rows *sql.Rows
	rows, err = db.Query("select h.guid,h.scheduled_time,h.execution_start_time,h.execution_end_time,h.message,h.state,h.schedule_guid,h.task_guid,h.created_at from histories h, schedulables a, schedules s, jobs j where h.schedule_guid=s.guid and s.schedulable_guid=a.guid and a.guid=j.guid and j.spaceguid like ? and j.name like ? order by h.execution_end_time desc", spaceguid, name)
	if err != nil {
		fmt.Printf("failed to query job histories: %s\n", err)
		return nil, err
	} else {
		result = histories2array(rows)
	}
	return result, nil
}

func GetCallHistories(spaceguid, name string) ([]model.History, error) {
	var err error
	if spaceguid == "" {
		spaceguid = "%"
	}
	if name == "" {
		name = "%"
	}
	result := make([]model.History, 0)
	db := GetDB()
	defer db.Close()
	var rows *sql.Rows
	rows, err = db.Query("select h.guid,h.scheduled_time,h.execution_start_time,h.execution_end_time,h.message,h.state,h.schedule_guid,h.task_guid,h.created_at from histories h, schedulables a, schedules s, calls c where h.schedule_guid=s.guid and s.schedulable_guid=a.guid and a.guid=c.guid and c.spaceguid=? and c.name=? order by h.execution_end_time desc", spaceguid, name)
	if err != nil {
		fmt.Printf("failed to query call histories: %s\n", err)
		return nil, err
	} else {
		result = histories2array(rows)
	}
	return result, nil
}

func histories2array(rows *sql.Rows) []model.History {
	result := make([]model.History, 0)
	if rows != nil {
		defer rows.Close()
		var guid, message, state, scheduleGuid, taskGuid string
		var scheduledTime, executionStartTime, executionEndTime, createdAt time.Time
		for rows.Next() {
			if err := rows.Scan(&guid, &scheduledTime, &executionStartTime, &executionEndTime, &message, &state, &scheduleGuid, &taskGuid, &createdAt); err != nil {
				fmt.Printf("failed to scan the history row:%s\n", err)
			} else {
				result = append(result, model.History{
					Guid:               guid,
					ScheduledTime:      scheduledTime,
					ExecutionStartTime: executionStartTime,
					ExecutionEndTime:   executionEndTime,
					Message:            message,
					State:              state,
					ScheduleGuid:       scheduleGuid,
					TaskGuid:           taskGuid,
					CreatedAt:          createdAt,
				})
			}
		}
	}
	return result
}

func DeleteHistoryByAge(maxDays int64) error {
	var err error
	nanoseconds := 1000000000 * 3600 * 24 * maxDays
	cutOffDate := time.Now().Add(time.Duration(-nanoseconds))
	fmt.Printf("cleaning histories table for rows older than %d days (older than %s)\n", maxDays, cutOffDate.Format(time.RFC3339))
	db := GetDB()
	defer db.Close()
	var result sql.Result
	if result, err = db.Exec("delete from histories where execution_start_time<?", cutOffDate); err != nil {
		fmt.Printf("failed to clean up histories: %s\n", err)
		return err
	}
	if numDeleted, _ := result.RowsAffected(); numDeleted > 0 {
		fmt.Printf("deleted %d rows from histories table\n", numDeleted)
	}
	return nil
}

func UpdateState(response model.TaskListResponse) error {
	var err error
	db := GetDB()
	defer db.Close()
	totErrCount := 0
	var totUpdates int64
	if response.Pagination.TotalResults > 0 {
		for _, resource := range response.Resources {
			message := ""
			if resource.State == "FAILED" {
				message = util.LastXChars(resource.Result.FailureReason, 255)
			}
			if result, err := db.Exec("update histories set state=?, message=? where task_guid=?", resource.State, message, resource.GUID); err != nil {
				totErrCount++
				fmt.Printf("failed to clean up histories: %s\n", err)
			} else {
				affectedRows, _ := result.RowsAffected()
				totUpdates = totUpdates + affectedRows
			}
		}
	}
	if totUpdates > 0 {
		fmt.Printf("updated %d rows in histories table\n", totUpdates)
	} else {
		fmt.Println()
	}
	if err != nil {
		return errors.New(fmt.Sprintf("%d updates failed, last failure reason:%s\n", totErrCount, err))
	}
	return err
}
