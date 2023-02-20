package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
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
	if serviceInstance, err := util.CfClient.GetServiceInstanceByGuid(serviceInstanceId); err != nil {
		fmt.Printf("could not find service instance with guid %s: %s\n", serviceInstanceId, err)
		util.WriteHttpResponse(w, http.StatusGone, model.DummyResponse{DummyField: "bye bye"})
	} else {
		if serviceBinding, err := util.CfClient.GetServiceBindingByGuid(serviceBindingId); err != nil {
			fmt.Printf("could not find service binding with guid %s: %s\n", serviceBindingId, err)
			util.WriteHttpResponse(w, http.StatusGone, model.DummyResponse{DummyField: "bye bye"})
		} else {
			// delete the jobs , calls and schedules for this app
			jobs, jobError := db.GetJobs(serviceInstance.SpaceGuid, "")
			_ = db.DeleteJobBySpaceGuidAndAppGuid(serviceInstance.SpaceGuid, serviceBinding.AppGuid)
			// delete the job from cron as well (we need the jobname to do that)
			var jobOrCallName string
			if jobError == nil {
				for _, job := range jobs {
					jobOrCallName = job.Name
					if job.AppGuid == serviceBinding.AppGuid {
						cron.DeleteJobByNameAndSpaceGuid(job.Name, serviceInstance.SpaceGuid)
					}
				}
			} else {
				fmt.Printf("no jobs found for deletion (app name: %s, spaceguid: %s)", jobOrCallName, serviceInstance.SpaceGuid)
			}

			calls, callError := db.GetCalls(serviceInstance.SpaceGuid, "")
			_ = db.DeleteCallBySpaceGuidAndAppGuid(serviceInstance.SpaceGuid, serviceBinding.AppGuid)
			// delete the call from cron as well (we need the callname to do that)
			if callError == nil {
				for _, call := range calls {
					jobOrCallName = call.Name
					if call.AppGuid == serviceBinding.AppGuid {
						cron.DeleteCallByNameAndSpaceGuid(call.Name, serviceInstance.SpaceGuid)
					}
				}
			} else {
				fmt.Printf("no calls found for deletion (app name: %s, spaceguid: %s)", jobOrCallName, serviceInstance.SpaceGuid)
			}
		}
	}
	fmt.Printf("delete service binding %s for service instance %s...\n", serviceBindingId, serviceInstanceId)
	util.WriteHttpResponse(w, http.StatusOK, model.DummyResponse{DummyField: "bye bye"})
}
