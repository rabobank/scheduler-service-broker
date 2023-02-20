package conf

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/rabobank/scheduler-service-broker/model"
	"os"
	"strconv"
)

var (
	Catalog                 model.Catalog
	ListenPort              int
	Debug                   = false
	httpTimeoutStr          = os.Getenv("SSB_HTTP_TIMEOUT")
	HttpTimeout             int
	HttpTimeoutDefault      = 10
	ClientId                = os.Getenv("SSB_CLIENT_ID")
	CfApiURL                = os.Getenv("SSB_CFAPI_URL")
	tokenRefreshIntervalStr = os.Getenv("SSB_TOKEN_REFRESH_INTERVAL")
	TokenRefreshInterval    int64
	SchedulerEndpoint       = os.Getenv("SSB_SCHEDULER_ENDPOINT")

	DebugStr      = os.Getenv("SSB_DEBUG")
	BrokerUser    = os.Getenv("SSB_BROKER_USER")
	CatalogDir    = os.Getenv("SSB_CATALOG_DIR")
	ListenPortStr = os.Getenv("SSB_LISTEN_PORT")

	DBUser              = os.Getenv("SSB_DB_USER")
	DBName              = os.Getenv("SSB_DB_NAME")
	DBHost              = os.Getenv("SSB_DB_HOST")
	maxHistoriesDaysStr = os.Getenv("SSB_MAX_HISTORY_DAYS")
	MaxHistoriesDays    int64

	BrokerPassword string
	DBPassword     string
	ClientSecret   string
)

const BasicAuthRealm = "scheduler-service-broker"

func EnvironmentComplete() {
	envComplete := true
	if DebugStr == "true" {
		Debug = true
	}
	if httpTimeoutStr == "" {
		HttpTimeout = HttpTimeoutDefault
	} else {
		var err error
		HttpTimeout, err = strconv.Atoi(httpTimeoutStr)
		if err != nil {
			fmt.Printf("failed reading envvar PCSB_HTTP_TIMEOUT, err: %s\n", err)
			envComplete = false
		}
	}
	if ClientId == "" {
		envComplete = false
		fmt.Println("missing envvar: SSB_CLIENT_ID")
	}
	if CfApiURL == "" {
		envComplete = false
		fmt.Println("missing envvar: SSB_CFAPI_URL")
	}
	if len(tokenRefreshIntervalStr) == 0 {
		TokenRefreshInterval = 90
	} else {
		var err error
		TokenRefreshInterval, err = strconv.ParseInt(tokenRefreshIntervalStr, 0, 64)
		if err != nil {
			panic(err)
		}
	}
	if SchedulerEndpoint == "" {
		envComplete = false
		fmt.Println("missing envvar: SSB_SCHEDULER_ENDPOINT")
	}
	if BrokerUser == "" {
		envComplete = false
		fmt.Println("missing envvar: SSB_BROKER_USER")
	}
	if CatalogDir == "" {
		CatalogDir = "catalog"
	}
	if ListenPortStr == "" {
		ListenPort = 8080
	} else {
		var err error
		ListenPort, err = strconv.Atoi(ListenPortStr)
		if err != nil {
			fmt.Printf("failed reading envvar LISTEN_PORT, err: %s\n", err)
			envComplete = false
		}
	}
	if DBUser == "" {
		envComplete = false
		fmt.Println("missing envvar: SSB_DB_USER")
	}
	if DBName == "" {
		DBName = "schedulerdb"
	}
	if DBHost == "" {
		DBHost = "localhost"
	}

	if len(maxHistoriesDaysStr) == 0 {
		MaxHistoriesDays = 90 // TODO make this a more reasonable default (100.000?)
	} else {
		var err error
		MaxHistoriesDays, err = strconv.ParseInt(maxHistoriesDaysStr, 0, 64)
		if err != nil {
			panic(err)
		}
	}

	if !envComplete {
		fmt.Println("one or more required environment variables missing, aborting...")
		os.Exit(8)
	}

	initCredentials()

}

// initCredentials - Get the credentials from credhub (VCAP_SERVICES envvar)
func initCredentials() {
	fmt.Println("getting credentials from credhub...")
	if appEnv, err := cfenv.Current(); err == nil {
		services, err := appEnv.Services.WithLabel("credhub")
		if err == nil {
			if len(services) != 1 {
				fmt.Printf("we expected exactly one bound credhub service instance, but found %d\n", len(services))
			} else {
				DBPassword = fmt.Sprint(services[0].Credentials["SSB_DB_PASSWORD"])
				BrokerPassword = fmt.Sprint(services[0].Credentials["SSB_BROKER_PASSWORD"])
				ClientSecret = fmt.Sprint(services[0].Credentials["SSB_CLIENT_SECRET"])
				allVarsFound := true
				if DBPassword == "" {
					fmt.Printf("credhub variable SSB_DB_PASSWORD is missing")
					allVarsFound = false
				}
				if BrokerPassword == "" {
					fmt.Printf("credhub variable SSB_BROKER_PASSWORD is missing")
					allVarsFound = false
				}
				if ClientSecret == "" {
					fmt.Printf("credhub variable SSB_CLIENT_SECRET is missing")
					allVarsFound = false
				}
				if !allVarsFound {
					os.Exit(8)
				}
			}
		} else {
			fmt.Printf("failed getting services from cf env: %s\n", err)
			os.Exit(8)
		}
	} else {
		fmt.Printf("failed to get the current cf env: %s\n", err)
		os.Exit(8)
	}
}
