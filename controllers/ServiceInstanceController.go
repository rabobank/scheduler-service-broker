package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rabobank/scheduler-service-broker/conf"
	"github.com/rabobank/scheduler-service-broker/model"
	"github.com/rabobank/scheduler-service-broker/util"
	"net/http"
	"strings"
)

const (
	ServiceInstanceStateSucceeded = "succeeded"
)

func Catalog(w http.ResponseWriter, r *http.Request) {
	_ = r // get rid of compiler warning
	util.WriteHttpResponse(w, http.StatusOK, conf.Catalog)
}

func GetServiceInstanceLastOperation(w http.ResponseWriter, _ *http.Request) {
	response := &model.LastOperation{
		State:       ServiceInstanceStateSucceeded,
		Description: "status is succeeded",
	}
	util.WriteHttpResponse(w, http.StatusOK, response)
}

func CreateServiceInstance(w http.ResponseWriter, r *http.Request) {
	serviceInstanceId := mux.Vars(r)["service_instance_guid"]
	fmt.Printf("create service instance for %s...\n", serviceInstanceId)
	var err error
	var serviceInstance model.ServiceInstance
	if err = util.ProvisionObjectFromRequest(r, &serviceInstance); err != nil {
		util.WriteHttpResponse(w, http.StatusBadRequest, model.BrokerError{Error: "FAILED", Description: err.Error(), InstanceUsable: false, UpdateRepeatable: false})
		return
	}
	serviceName := util.GetServiceById(serviceInstance.ServiceId).Name
	if !strings.HasPrefix(serviceName, "scheduler") {
		util.WriteHttpResponse(w, http.StatusBadRequest, model.BrokerError{Error: "FAILED", Description: fmt.Sprintf("service %s is not supported", serviceName), InstanceUsable: false, UpdateRepeatable: false})
		return
	}
	lastOperation := &model.LastOperation{State: "in progress", Description: "service instance create is in progress ..."}
	// providing the dashboardURL is actually the only thing that we need to implement for the broker API
	response := model.CreateServiceInstanceResponse{LastOperation: lastOperation, DashboardUrl: conf.SchedulerEndpoint}
	util.WriteHttpResponse(w, http.StatusCreated, response)
	return
}

func DeleteServiceInstance(w http.ResponseWriter, r *http.Request) {
	serviceInstanceId := mux.Vars(r)["service_instance_guid"]
	fmt.Printf("delete service instance %s...\n", serviceInstanceId)
	if serviceInstanceId == "not exist" {
		util.WriteHttpResponse(w, http.StatusGone, model.DeleteServiceInstanceResponse{Result: fmt.Sprintf("service instance with guid %s not found", serviceInstanceId)})
		return
	}

	// there is nothing for us to do here
	util.WriteHttpResponse(w, http.StatusOK, model.DeleteServiceInstanceResponse{})
}
