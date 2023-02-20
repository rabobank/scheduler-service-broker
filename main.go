package main

import (
	"encoding/json"
	"fmt"
	"github.com/rabobank/scheduler-service-broker/conf"
	"github.com/rabobank/scheduler-service-broker/cron"
	"github.com/rabobank/scheduler-service-broker/db"
	"github.com/rabobank/scheduler-service-broker/server"
	"github.com/rabobank/scheduler-service-broker/util"
	"os"
)

func main() {
	fmt.Printf("scheduler-service-broker starting, version:%s, commit:%s\n", conf.VERSION, conf.COMMIT)

	conf.EnvironmentComplete()

	initialize()

	cron.StartRunner()

	cron.StartHousekeeping()

	server.StartServer()
}

// initialize scheduler-service-broker:
func initialize() {
	catalogFile := fmt.Sprintf("%s/catalog.json", conf.CatalogDir)
	if file, err := os.ReadFile(catalogFile); err != nil {
		fmt.Printf("failed reading catalog file %s: %s\n", catalogFile, err)
		os.Exit(8)
	} else {
		if err = json.Unmarshal(file, &conf.Catalog); err != nil {
			fmt.Printf("failed unmarshalling catalog file %s, error: %s\n", catalogFile, err)
			os.Exit(8)
		} else {

			// login to cf and get a client handle
			util.CfClient = *util.GetCFClient()

			// test if the DB can be opened
			database := db.GetDB()
			if err = database.Close(); err != nil {
				fmt.Printf("failed to open database %s:%s\n", fmt.Sprintf("%s:%s@(%s)/%s?parseTime=true", conf.DBUser, "<redacted>", conf.DBHost, conf.DBName), err)
				os.Exit(8)
			}
		}
	}
}
