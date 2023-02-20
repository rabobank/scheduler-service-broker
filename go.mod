module github.com/rabobank/scheduler-service-broker

go 1.20

require (
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20211117203709-9b81b3940cc7
	github.com/cloudfoundry-community/go-cfenv v1.18.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/gorilla/context v1.1.1
	github.com/gorilla/mux v1.8.0
)

require (
	code.cloudfoundry.org/gofileutils v0.0.0-20170111115228-4d0c80011a0f // indirect
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/golang/protobuf v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/oauth2 v0.0.0-20190130055435-99b60b757ec1 // indirect
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/appengine v1.4.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

exclude (
	golang.org/x/text v0.3.0
	golang.org/x/text v0.3.3
	gopkg.in/yaml.v2 v2.2.1
)
