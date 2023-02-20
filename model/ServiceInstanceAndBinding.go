package model

type ServiceInstance struct {
	ServiceId        string      `json:"service_id"`
	PlanId           string      `json:"plan_id"`
	OrganizationGuid string      `json:"organization_guid"`
	SpaceGuid        string      `json:"space_guid"`
	Context          *Context    `json:"context"`
	Parameters       *Parameters `json:"parameters,omitempty"`
}

type CreateServiceInstanceResponse struct {
	ServiceId     string         `json:"service_id"`
	PlanId        string         `json:"plan_id"`
	DashboardUrl  string         `json:"dashboard_url"`
	LastOperation *LastOperation `json:"last_operation,omitempty"`
}

type DeleteServiceInstanceResponse struct {
	Result string `json:"result,omitempty"`
}

// Parameters These are the potential parameters that can be given on the -c parameter of "cf create-service"
type Parameters struct {
	// there are no parameters for this broker
}

type CreateServiceBindingResponse struct {
	// SyslogDrainUrl string      `json:"syslog_drain_url, omitempty"`
	Credentials *Credentials `json:"credentials"`
}

type Credentials struct {
	DummyCreds string `json:"dummycreds"`
}

type DummyResponse struct {
	DummyField string
}
