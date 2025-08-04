package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/rabobank/scheduler-service-broker/util"
)

type Call struct {
	Guid        string
	AppGuid     string
	SpaceGuid   string
	State       string
	Name        string
	Url         string
	AuthHeader  string
	Schedulable int64
}

func (call Call) String() string {
	return fmt.Sprintf("Guid:%s AppGuid:%s, SpaceGuid:%s, State:%s, Name:%s, Url:%s, AuthHeader: <redacted>", call.Guid, call.AppGuid, call.SpaceGuid, call.State, call.Name, call.Url)
}

func InsertCall(call Call) (string, error) {
	var err error
	db := GetDB()
	defer func() { _ = db.Close() }()
	newGuid := util.GenerateGUID()
	if _, err = db.Exec("insert into schedulables(guid) values(?)", newGuid); err != nil {
		fmt.Printf("failed to insert schedulable, error: %s\n", err)
	} else {
		if _, err = db.Exec("insert into calls(guid,appguid,spaceguid,name,url,authheader) values(?,?,?,?,?,?)", newGuid, call.AppGuid, call.SpaceGuid, call.Name, call.Url, call.AuthHeader); err != nil {
			fmt.Printf("failed to insert %v, error: %s\n", call, err)
			_, _ = db.Exec("delete from schedulables where guid=?", newGuid)
		} else {
			call.Guid = newGuid
			fmt.Printf("inserted %v\n", call)
		}
	}
	return call.Guid, err
}

func GetCalls(spaceguid, name string) ([]Call, error) {
	var err error
	result := make([]Call, 0)
	if spaceguid == "" {
		spaceguid = "%"
	}
	if name == "" {
		name = "%"
	}
	db := GetDB()
	defer func() { _ = db.Close() }()
	var rows *sql.Rows
	rows, err = db.Query("select guid,appguid,spaceguid,state,name,url,authheader from calls where spaceguid like ? and name like ?", spaceguid, name)
	if err != nil {
		fmt.Printf("failed to query the calls, err: %s\n", err)
		return nil, err
	} else {
		result = calls2array(rows)
	}
	return result, nil
}

func calls2array(rows *sql.Rows) []Call {
	result := make([]Call, 0)
	if rows != nil {
		defer func() { _ = rows.Close() }()
		var guid, appguid, spaceguid, state, name, url, authheader string
		for rows.Next() {
			err := rows.Scan(&guid, &appguid, &spaceguid, &state, &name, &url, &authheader)
			if err != nil {
				fmt.Printf("failed to scan the call row, error:%s\n", err)
			} else {
				result = append(result, Call{
					Guid:       guid,
					AppGuid:    appguid,
					SpaceGuid:  spaceguid,
					State:      state,
					Name:       name,
					Url:        url,
					AuthHeader: authheader,
				})
			}
		}
	}
	return result
}

func DeleteCallBySpaceGuidAndCallname(spaceguid, callname string) error {
	var err error
	db := GetDB()
	defer func() { _ = db.Close() }()
	// delete the schedulable, call will be cascade-deleted, if there are still schedules that "run" this call, they will also be cascade-deleted
	result, err := db.Exec("delete from schedulables where guid in (select guid from calls where name=? and spaceguid=?)", callname, spaceguid)
	numDeletes, _ := result.RowsAffected()
	if numDeletes == 0 {
		err = errors.New(fmt.Sprintf("call %s does not exist, no rows deleted", callname))
	}
	if err != nil {
		fmt.Printf("failed to delete call/schedules for callname %s and spaceguid %s, numDeletes=%d error: %s\n", callname, spaceguid, numDeletes, err)
		return err
	} else {
		fmt.Printf("deleted %d call/schedule for callname %s and spaceguid %s\n", numDeletes, callname, spaceguid)
		return nil
	}
}

func DeleteCallBySpaceGuidAndAppGuid(spaceguid, appguid string) error {
	var err error
	db := GetDB()
	defer func() { _ = db.Close() }()
	// delete the schedulable, call will be cascade-deleted, if there are still schedules that "run" this call, they will also be cascade-deleted
	result, err := db.Exec("delete from schedulables where guid in (select guid from calls where appguid=? and spaceguid=?)", appguid, spaceguid)
	numDeletes, _ := result.RowsAffected()
	if numDeletes == 0 {
		err = errors.New(fmt.Sprintf("call for appguid %s and spaceguid %s does not exist, no rows deleted", appguid, spaceguid))
	}
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	} else {
		fmt.Printf("deleted %d call/schedules for appguid %s and spaceguid %s\n", numDeletes, appguid, spaceguid)
		return nil
	}
}
