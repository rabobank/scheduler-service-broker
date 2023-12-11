package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/rabobank/scheduler-service-broker/util"
)

type Job struct {
	Guid       string
	AppGuid    string
	SpaceGuid  string
	State      string
	Name       string
	Command    string
	MemoryInMB int
	DiskInMB   int
}

func (job Job) String() string {
	return fmt.Sprintf("Guid:%s AppGuid:%s, SpaceGuid:%s, State:%s, Name:%s, Command:%s, MemoryInMB: %d, DiskInMB: %d", job.Guid, job.AppGuid, job.SpaceGuid, job.State, job.Name, job.Command, job.MemoryInMB, job.DiskInMB)
}

func InsertJob(job Job) (string, error) {
	var err error
	db := GetDB()
	defer func() { _ = db.Close() }()
	newGuid := util.GenerateGUID()
	if _, err = db.Exec("insert into schedulables(guid) values(?)", newGuid); err != nil {
		fmt.Printf("failed to insert schedulable, error: %s\n", err)
	} else {
		if _, err = db.Exec("insert into jobs(guid,appguid,spaceguid,name,command,memoryinmb,diskinmb) values(?,?,?,?,?,?,?)", newGuid, job.AppGuid, job.SpaceGuid, job.Name, job.Command, job.MemoryInMB, job.DiskInMB); err != nil {
			fmt.Printf("failed to insert %v, error: %s\n", job, err)
			_, _ = db.Exec("delete from schedulables where guid=?", newGuid)
		} else {
			job.Guid = newGuid
			fmt.Printf("inserted %v\n", job)
		}
	}
	return job.Guid, err
}

func GetJobs(spaceguid, name string) ([]Job, error) {
	var err error
	result := make([]Job, 0)
	if spaceguid == "" {
		spaceguid = "%"
	}
	if name == "" {
		name = "%"
	}
	db := GetDB()
	defer func() { _ = db.Close() }()
	var rows *sql.Rows
	rows, err = db.Query("select guid,appguid,spaceguid,state,name,command,memoryinmb,diskinmb from jobs where spaceguid like ? and name like ?", spaceguid, name)
	if err != nil {
		fmt.Printf("failed to query the jobs, err: %s\n", err)
		return nil, err
	} else {
		result = jobs2array(rows)
	}
	return result, nil
}

func jobs2array(rows *sql.Rows) []Job {
	result := make([]Job, 0)
	if rows != nil {
		defer func() { _ = rows.Close() }()
		var guid, appguid, spaceguid, state, name, command string
		var memoryinmb, diskinmb int
		for rows.Next() {
			if err := rows.Scan(&guid, &appguid, &spaceguid, &state, &name, &command, &memoryinmb, &diskinmb); err != nil {
				fmt.Printf("failed to scan the job row, error:%s\n", err)
			} else {
				result = append(result, Job{
					Guid:       guid,
					AppGuid:    appguid,
					SpaceGuid:  spaceguid,
					State:      state,
					Name:       name,
					Command:    command,
					MemoryInMB: memoryinmb,
					DiskInMB:   diskinmb,
				})
			}
		}
	}
	return result
}

func DeleteJobBySpaceGuidAndJobname(spaceguid, jobname string) error {
	var err error
	db := GetDB()
	defer func() { _ = db.Close() }()
	// delete the schedulable, job will be cascade-deleted, if there are still schedules that "run" this job, they will also be cascade-deleted
	result, err := db.Exec("delete from schedulables where guid in (select guid from jobs where name=? and spaceguid=?)", jobname, spaceguid)
	numDeletes, _ := result.RowsAffected()
	if numDeletes == 0 {
		err = errors.New(fmt.Sprintf("job %s does not exist, no rows deleted", jobname))
	}
	if err != nil {
		fmt.Printf("failed to delete job/schedules for jobname %s and spaceguid %s, error: %s\n", jobname, spaceguid, err)
		return err
	}
	fmt.Printf("deleted %d job/schedule for jobname %s and spaceguid %s\n", numDeletes, jobname, spaceguid)
	return nil
}

func DeleteJobBySpaceGuidAndAppGuid(spaceguid, appguid string) error {
	var err error
	db := GetDB()
	defer func() { _ = db.Close() }()
	// delete the schedulable, job will be cascade-deleted, if there are still schedules that "run" this job, they will also be cascade-deleted
	result, err := db.Exec("delete from schedulables where guid in (select guid from jobs where appguid=? and spaceguid=?)", appguid, spaceguid)
	numDeletes, _ := result.RowsAffected()
	if numDeletes == 0 {
		err = errors.New(fmt.Sprintf("job for appguid %s and spaceguid %s does not exist, no rows deleted", appguid, spaceguid))
	}
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}
	fmt.Printf("deleted %d job/schedules for appguid %s and spaceguid %s\n", numDeletes, appguid, spaceguid)
	return nil
}
