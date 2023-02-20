package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rabobank/scheduler-service-broker/conf"
	"time"
)

func GetDB() (db *sql.DB) {
	var err error
	dbDriver := "mysql"
	if db, err = sql.Open(dbDriver, fmt.Sprintf("%s:%s@(%s)/%s?parseTime=true", conf.DBUser, conf.DBPassword, conf.DBHost, conf.DBName)); err != nil {
		panic(err.Error())
	} else {
		return db
	}
}

type Schedulable struct {
	Id        int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
