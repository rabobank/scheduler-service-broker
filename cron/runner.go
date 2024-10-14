package cron

import (
	"fmt"
	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/cloudfoundry/go-cfclient/v3/resource"
	"github.com/rabobank/scheduler-service-broker/conf"
	"github.com/rabobank/scheduler-service-broker/db"
	"github.com/rabobank/scheduler-service-broker/model"
	"github.com/rabobank/scheduler-service-broker/util"
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var jobRunner *cron.Cron
var callRunner *cron.Cron
var jobRunnerEntries map[cron.EntryID]string  // value is a concatenation of spaceguid, scheduleguid and job/call name
var callRunnerEntries map[cron.EntryID]string // value is a concatenation of spaceguid, scheduleguid and job/call name

type Job struct {
	sjob model.SchedulableJob
}

// Run - This function is required for Job and will be called by cron.Cron
func (job Job) Run() {
	var scheduledTime time.Time
	// find the scheduled time from list of cron entries, and use that to be inserted later on in the histories table
	for _, entry := range jobRunner.Entries() {
		if entry.Job.(Job).sjob.ScheduleGuid == job.sjob.ScheduleGuid {
			scheduledTime = entry.Prev
		}
	}
	DoJob(scheduledTime, job.sjob)
}

type Call struct {
	scall model.SchedulableCall
}

// Run - This function is required for Call and will be called by cron.Cron
func (call Call) Run() {
	var scheduledTime time.Time
	// find the scheduled time from list of cron entries, and use that to be inserted later on in the histories table
	for _, entry := range callRunner.Entries() {
		if entry.Job.(Call).scall.ScheduleGuid == call.scall.ScheduleGuid {
			scheduledTime = entry.Prev
		}
	}
	DoCall(scheduledTime, call.scall)
}

// AddJob - Add a Runnable job to the list of cron entries
func AddJob(sjob model.SchedulableJob) {
	job := Job{sjob}
	if entryID, err := jobRunner.AddJob(sjob.Expression, job); err != nil {
		fmt.Printf("failed to add func with schedule %s, jobname %s, appguid %s, spaceguid %s: %s\n", job.sjob.Expression, job.sjob.JobName, job.sjob.AppGuid, job.sjob.SpaceGuid, err)
	} else {
		key := fmt.Sprintf("%s-%s-%s", sjob.ScheduleGuid, sjob.SpaceGuid, sjob.JobName)
		jobRunnerEntries[entryID] = key
		fmt.Printf("job (id=%d [key=%s]) with jobname %s (mem=%d, disk=%d) added with schedule %s\n", entryID, key, sjob.JobName, sjob.MemoryInMB, sjob.DiskInMB, sjob.Expression)
	}
}

// DeleteJobByNameAndSpaceGuid - Delete a Runnable job from the list of cron Entries
func DeleteJobByNameAndSpaceGuid(jobname, spaceGuid string) {
	for k, v := range jobRunnerEntries {
		if strings.Contains(v, spaceGuid) && strings.Contains(v, jobname) {
			jobRunner.Remove(k)
			delete(jobRunnerEntries, k)
			fmt.Printf("removed job schedule (id=%d) for spaceGuid %s, %d schedules left\n", k, spaceGuid, len(jobRunner.Entries()))
		}
	}
}

// DeleteJobByScheduleGuid - Delete a Runnable job from the list of cron Entries
func DeleteJobByScheduleGuid(scheduleGuid string) {
	for k, v := range jobRunnerEntries {
		if strings.Contains(v, scheduleGuid) {
			jobRunner.Remove(k)
			delete(jobRunnerEntries, k)
			fmt.Printf("removed job schedule (id=%d) for scheduleGuid %s, %d schedules left\n", k, scheduleGuid, len(jobRunner.Entries()))
		}
	}
}

// AddCall - Add a http call to the list of cron entries
func AddCall(scall model.SchedulableCall) {
	call := Call{scall}
	if entryID, err := callRunner.AddJob(scall.Expression, call); err != nil {
		fmt.Printf("failed to add func with schedule %s, jobname %s, appguid %s, spaceguid %s: %s\n", call.scall.Expression, call.scall.CallName, call.scall.AppGuid, call.scall.SpaceGuid, err)
	} else {
		key := fmt.Sprintf("%s-%s-%s", scall.ScheduleGuid, scall.SpaceGuid, scall.CallName)
		callRunnerEntries[entryID] = key
		fmt.Printf("call (id=%d [key=%s]) with callname %s added with schedule %s, number of entries: %d/%d\n", entryID, key, scall.CallName, scall.Expression, len(callRunner.Entries()), len(callRunnerEntries))
	}
}

// DeleteCallByNameAndSpaceGuid - Delete a Runnable job from the list of cron Entries
func DeleteCallByNameAndSpaceGuid(callname, spaceGuid string) {
	for k, v := range callRunnerEntries {
		if strings.Contains(v, spaceGuid) && strings.Contains(v, callname) {
			callRunner.Remove(k)
			delete(callRunnerEntries, k)
			fmt.Printf("removed call schedule (id=%d) for spaceGuid %s, schedules left: %d/%d\n", k, spaceGuid, len(callRunner.Entries()), len(callRunnerEntries))
		}
	}
}

// DeleteCallByScheduleGuid - Delete a Runnable call from the list of cron Entries
func DeleteCallByScheduleGuid(scheduleGuid string) {
	for k, v := range callRunnerEntries {
		if strings.Contains(v, scheduleGuid) {
			callRunner.Remove(k)
			delete(callRunnerEntries, k)
			fmt.Printf("removed call schedule (id=%d) for scheduleGuid %s, schedules left: %d/%d\n", k, scheduleGuid, len(callRunner.Entries()), len(callRunnerEntries))
		}
	}
}

// StartRunner - The main function to start the cron job. It will start reading the database for schedules and add them all to cron.
func StartRunner() {
	jobRunnerEntries = make(map[cron.EntryID]string)
	callRunnerEntries = make(map[cron.EntryID]string)
	if schedulableJobs, err := db.GetSchedulableJobs(); err != nil {
		panic(fmt.Sprintf("failed to get schedulable jobs: %s", err))
	} else {
		if schedulableCalls, err := db.GetSchedulableCalls(); err != nil {
			panic(fmt.Sprintf("failed to get schedulable calls: %s", err))
		} else {
			if conf.Debug {
				jobRunner = cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "jobRunner: ", log.Ldate|log.Ltime|log.LUTC))), cron.WithLocation(time.UTC))
				callRunner = cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "callRunner: ", log.Ldate|log.Ltime|log.LUTC))), cron.WithLocation(time.UTC))
			} else {
				jobRunner = cron.New(cron.WithLocation(time.UTC))
				callRunner = cron.New(cron.WithLocation(time.UTC))
			}
			for _, schedulableJob := range schedulableJobs {
				job := Job{schedulableJob}
				if entryId, err := jobRunner.AddJob(job.sjob.Expression, job); err != nil {
					fmt.Printf("failed to add func with schedule %s, jobname %s, appguid %s, spaceguid %s: %s\n", job.sjob.Expression, job.sjob.JobName, job.sjob.AppGuid, job.sjob.SpaceGuid, err)
				} else {
					jobRunnerEntries[entryId] = fmt.Sprintf("%s-%s-%s", schedulableJob.ScheduleGuid, schedulableJob.SpaceGuid, schedulableJob.JobName)
					fmt.Printf("job %s added with schedule(%d) %s\n", job.sjob.JobName, entryId, job.sjob.Expression)
				}
			}
			for _, schedulableCall := range schedulableCalls {
				call := Call{schedulableCall}
				if entryId, err := callRunner.AddJob(call.scall.Expression, call); err != nil {
					fmt.Printf("failed to add func with schedule %s, jobname %s, appguid %s, spaceguid %s: %s\n", call.scall.Expression, call.scall.CallName, call.scall.AppGuid, call.scall.SpaceGuid, err)
				} else {
					callRunnerEntries[entryId] = fmt.Sprintf("%s-%s-%s", schedulableCall.ScheduleGuid, schedulableCall.SpaceGuid, schedulableCall.CallName)
					fmt.Printf("call %s added with schedule(%d) %s\n", call.scall.CallName, entryId, call.scall.Expression)
				}
			}

			jobRunner.Start()
			fmt.Println("JobRunner (re)started")
			callRunner.Start()
			fmt.Println("CallRunner (re)started")
		}
	}
}

// StartHousekeeping - This function starts a loop in his own routine to query cf tasks, and update the histories table with the task state and message, and to clean up histories and schedules table.
func StartHousekeeping() {
	var loopCounter int
	go func() {
		for {
			time.Sleep(time.Minute * 1)
			loopCounter++

			util.PrintIfDebug("copying cf task states to histories...\n")
			if loopCounter%10 == 0 {
				// once every 10 loops, clean up histories table
				_ = db.DeleteHistoryByAge(conf.MaxHistoriesDays)
				_ = db.DeleteScheduleByAge(conf.MaxHistoriesDays)
			}
			taskListOptions := client.TaskListOptions{ListOptions: &client.ListOptions{PerPage: 100, OrderBy: "-created_at"}}
			if tasks, err := conf.CfClient.Tasks.ListAll(conf.CfCtx, &taskListOptions); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("found %d cf tasks", len(tasks))
				if err = db.UpdateState(tasks); err != nil {
					fmt.Printf("failed to update histories state: %s\n", err)
				}
			}
		}
	}()
}

func DoCall(scheduledTime time.Time, call model.SchedulableCall) {
	historyRecord := model.History{
		Guid:               util.GenerateGUID(),
		ScheduledTime:      scheduledTime,
		ExecutionStartTime: time.Now(),
		ScheduleGuid:       call.ScheduleGuid,
		CreatedAt:          time.Now(),
	}
	// do the actual call to the URL
	transport := http.Transport{IdleConnTimeout: time.Second}
	httpClient := http.Client{Timeout: time.Duration(conf.HttpTimeout) * time.Second, Transport: &transport}
	callUrl, _ := url.Parse(call.Url)
	req := http.Request{Method: http.MethodGet, URL: callUrl}
	if call.AuthHeader != "" {
		req = http.Request{Method: http.MethodGet, URL: callUrl, Header: map[string][]string{"Authorization": {call.AuthHeader}}}
	}
	resp, err := httpClient.Do(&req)
	historyRecord.ExecutionEndTime = time.Now()
	if err != nil {
		fmt.Printf("failed calling url \"%s\": %s\n", callUrl, err)
		historyRecord.State = "FAILED"
		historyRecord.Message = err.Error()
	}
	if resp != nil {
		body, _ := io.ReadAll(resp.Body)
		util.PrintfIfDebug("response from call: %s\n", util.LastXChars(fmt.Sprintf("%s", body), 1024))
		historyRecord.State = resp.Status
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("failed calling url \"%s\": %s\n", callUrl, resp.Status)
			historyRecord.Message = util.LastXChars(fmt.Sprintf("%s", body), 255)
		} else {
			historyRecord.Message = ""
		}
	}
	if call.ScheduleGuid == "" {
		// only when we did a "cf run-call"
		if historyRecord.ScheduleGuid, _, err = db.InsertCallSchedule(model.ScheduleRequest{SpaceGUID: call.SpaceGuid, Name: call.CallName, CronExpression: "none", ExpressionType: "execute"}); err != nil {
			fmt.Printf("failed to insert execute-type call schedule: %s\n", err)
		}
	}
	_, _ = db.InsertHistory(historyRecord)
	fmt.Printf("call to url \"%s\" executed\n", callUrl)
}

func DoJob(scheduledTime time.Time, job model.SchedulableJob) {
	historyRecord := model.History{Guid: util.GenerateGUID(), ScheduledTime: scheduledTime, ExecutionStartTime: time.Now(), ScheduleGuid: job.ScheduleGuid, CreatedAt: time.Now()}

	// run the actual cf task
	taskCreateRequest := resource.TaskCreate{Command: &job.Command, MemoryInMB: &job.MemoryInMB, DiskInMB: &job.DiskInMB}
	if task, err := conf.CfClient.Tasks.Create(conf.CfCtx, job.AppGuid, &taskCreateRequest); err != nil {
		fmt.Printf("failed running cmd %s in app with guid %s: %s\n", job.Command, job.AppGuid, err)
	} else {
		historyRecord.ExecutionEndTime = time.Now()
		historyRecord.State = task.State
		historyRecord.TaskGuid = task.GUID
		historyRecord.Message = ""
		if task.Result.FailureReason != nil {
			historyRecord.Message = util.LastXChars(fmt.Sprintf("%s", *task.Result.FailureReason), 255)
		}
	}

	if job.ScheduleGuid == "" {
		// only when we did a "cf run-job"
		var err error
		if historyRecord.ScheduleGuid, _, err = db.InsertJobSchedule(model.ScheduleRequest{SpaceGUID: job.SpaceGuid, Name: job.JobName, CronExpression: "none", ExpressionType: "execute"}); err != nil {
			fmt.Printf("failed to insert execute-type job schedule: %s\n", err)
		}
	}
	_, _ = db.InsertHistory(historyRecord)
	fmt.Printf("cf task with cmd \"%s\" started in app with guid %s\n", job.Command, job.AppGuid)
}
