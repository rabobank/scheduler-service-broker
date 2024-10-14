package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rabobank/scheduler-service-broker/conf"
	"github.com/rabobank/scheduler-service-broker/cron"
	"github.com/rabobank/scheduler-service-broker/db"
	"github.com/rabobank/scheduler-service-broker/model"
	"github.com/rabobank/scheduler-service-broker/util"
	"net/http"
)

func GetServiceBinding(w http.ResponseWriter, r *http.Request) {
	serviceInstanceId := mux.Vars(r)["service_instance_guid"]
	serviceBindingId := mux.Vars(r)["service_binding_guid"]
	fmt.Printf("get service binding %s for service instance %s...\n", serviceBindingId, serviceInstanceId)
	util.WriteHttpResponse(w, http.StatusOK, model.DummyResponse{DummyField: "bye bye"})
}

func CreateServiceBinding(w http.ResponseWriter, r *http.Request) {
	serviceInstanceId := mux.Vars(r)["service_instance_guid"]
	serviceBindingId := mux.Vars(r)["service_binding_guid"]
	fmt.Printf("create service binding %s for service instance %s...\n", serviceBindingId, serviceInstanceId)
	util.WriteHttpResponse(w, http.StatusOK, model.CreateServiceBindingResponse{Credentials: nil})
}

func DeleteServiceBinding(w http.ResponseWriter, r *http.Request) {
	serviceInstanceId := mux.Vars(r)["service_instance_guid"]
	serviceBindingId := mux.Vars(r)["service_binding_guid"]
	if serviceInstance, err := conf.CfClient.ServiceInstances.Get(conf.CfCtx, serviceInstanceId); err != nil {
		fmt.Printf("could not find service instance with guid %s: %s\n", serviceInstanceId, err)
		util.WriteHttpResponse(w, http.StatusGone, model.DummyResponse{DummyField: "bye bye"})
	} else {
		if serviceBinding, err := conf.CfClient.ServiceCredentialBindings.Get(conf.CfCtx, serviceBindingId); err != nil {
			fmt.Printf("could not find service binding with guid %s: %s\n", serviceBindingId, err)
			util.WriteHttpResponse(w, http.StatusGone, model.DummyResponse{DummyField: "bye bye"})
		} else {
			// delete the jobs , calls and schedules for this app
			jobs, jobError := db.GetJobs(serviceInstance.Relationships.Space.Data.GUID, "")
			_ = db.DeleteJobBySpaceGuidAndAppGuid(serviceInstance.Relationships.Space.Data.GUID, serviceBinding.Relationships.App.Data.GUID)
			// delete the job from cron as well (we need the jobname to do that)
			var jobOrCallName string
			if jobError == nil {
				for _, job := range jobs {
					jobOrCallName = job.Name
					if job.AppGuid == serviceBinding.Relationships.App.Data.GUID {
						cron.DeleteJobByNameAndSpaceGuid(job.Name, serviceInstance.Relationships.Space.Data.GUID)
					}
				}
			} else {
				fmt.Printf("no jobs found for deletion (app name: %s, spaceguid: %s)", jobOrCallName, serviceInstance.Relationships.Space.Data.GUID)
			}

			calls, callError := db.GetCalls(serviceInstance.Relationships.Space.Data.GUID, "")
			_ = db.DeleteCallBySpaceGuidAndAppGuid(serviceInstance.Relationships.Space.Data.GUID, serviceBinding.Relationships.App.Data.GUID)
			// delete the call from cron as well (we need the callname to do that)
			if callError == nil {
				for _, call := range calls {
					jobOrCallName = call.Name
					if call.AppGuid == serviceBinding.Relationships.App.Data.GUID {
						cron.DeleteCallByNameAndSpaceGuid(call.Name, serviceInstance.Relationships.Space.Data.GUID)
					}
				}
			} else {
				fmt.Printf("no calls found for deletion (app name: %s, spaceguid: %s)", jobOrCallName, serviceInstance.Relationships.Space.Data.GUID)
			}
		}
	}
	fmt.Printf("delete service binding %s for service instance %s...\n", serviceBindingId, serviceInstanceId)
	util.WriteHttpResponse(w, http.StatusOK, model.DummyResponse{DummyField: "bye bye"})
}
