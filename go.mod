module github.com/rabobank/scheduler-service-broker

go 1.23
replace (
	golang.org/x/text => golang.org/x/text v0.22.0
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/cloudfoundry-community/go-cfenv v1.18.0
	github.com/cloudfoundry/go-cfclient/v3 v3.0.0-alpha.10
	github.com/go-sql-driver/mysql v1.9.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/gorilla/context v1.1.2
	github.com/gorilla/mux v1.8.1
	github.com/robfig/cron/v3 v3.0.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/codegangsta/inject v0.0.0-20150114235600-33e0aa1cb7c0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-martini/martini v0.0.0-20170121215854-22fa46961aab // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/martini-contrib/render v0.0.0-20150707142108-ec18f8345a11 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/oauth2 v0.27.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
