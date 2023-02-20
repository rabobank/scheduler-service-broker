package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rabobank/scheduler-service-broker/conf"
	"github.com/rabobank/scheduler-service-broker/controllers"
	"net/http"
	"os"
)

// StartServer - We start an httpd with 2 routers, one for the cf service-broker api's, that one uses Basic Authentication, and another one for all the CRUD requests for the scheduler, that one requires a valid JWT.
func StartServer() {
	brokerRouter := mux.NewRouter()
	brokerRouter.Use(controllers.DebugMiddleware)
	brokerRouter.Use(controllers.BasicAuthMiddleware)
	brokerRouter.Use(controllers.AddHeadersMiddleware)
	brokerRouter.HandleFunc("/v2/catalog", controllers.Catalog).Methods(http.MethodGet)
	brokerRouter.HandleFunc("/v2/service_instances/{service_instance_guid}/last_operation", controllers.GetServiceInstanceLastOperation).Methods(http.MethodGet)
	brokerRouter.HandleFunc("/v2/service_instances/{service_instance_guid}", controllers.CreateServiceInstance).Methods(http.MethodPut)
	brokerRouter.HandleFunc("/v2/service_instances/{service_instance_guid}", controllers.DeleteServiceInstance).Methods(http.MethodDelete)
	brokerRouter.HandleFunc("/v2/service_instances/{service_instance_guid}/service_bindings/{service_binding_guid}", controllers.GetServiceBinding).Methods(http.MethodGet)
	brokerRouter.HandleFunc("/v2/service_instances/{service_instance_guid}/service_bindings/{service_binding_guid}", controllers.CreateServiceBinding).Methods(http.MethodPut)
	brokerRouter.HandleFunc("/v2/service_instances/{service_instance_guid}/service_bindings/{service_binding_guid}", controllers.DeleteServiceBinding).Methods(http.MethodDelete)
	http.Handle("/v2/", brokerRouter)

	apiRouter := mux.NewRouter()
	apiRouter.Use(controllers.DebugMiddleware)
	apiRouter.Use(controllers.CheckJWTMiddleware)
	apiRouter.Use(controllers.AddHeadersMiddleware)
	apiRouter.HandleFunc("/api/jobs", controllers.JobGet).Methods(http.MethodGet)
	apiRouter.HandleFunc("/api/jobs", controllers.JobCreate).Methods(http.MethodPost)
	apiRouter.HandleFunc("/api/jobs", controllers.JobDelete).Methods(http.MethodDelete)
	apiRouter.HandleFunc("/api/jobs", controllers.JobRun).Methods(http.MethodPut)
	apiRouter.HandleFunc("/api/calls", controllers.CallGet).Methods(http.MethodGet)
	apiRouter.HandleFunc("/api/calls", controllers.CallCreate).Methods(http.MethodPost)
	apiRouter.HandleFunc("/api/calls", controllers.CallDelete).Methods(http.MethodDelete)
	apiRouter.HandleFunc("/api/calls", controllers.CallRun).Methods(http.MethodPut)
	apiRouter.HandleFunc("/api/jobschedules", controllers.JobScheduleGet).Methods(http.MethodGet)
	apiRouter.HandleFunc("/api/callschedules", controllers.CallScheduleGet).Methods(http.MethodGet)
	apiRouter.HandleFunc("/api/jobschedules", controllers.JobScheduleCreate).Methods(http.MethodPost)
	apiRouter.HandleFunc("/api/callschedules", controllers.CallScheduleCreate).Methods(http.MethodPost)
	apiRouter.HandleFunc("/api/jobschedules", controllers.JobScheduleDelete).Methods(http.MethodDelete)
	apiRouter.HandleFunc("/api/callschedules", controllers.CallScheduleDelete).Methods(http.MethodDelete)
	apiRouter.HandleFunc("/api/jobhistories", controllers.JobHistoriesGet).Methods(http.MethodGet)
	apiRouter.HandleFunc("/api/callhistories", controllers.CallHistoriesGet).Methods(http.MethodGet)
	http.Handle("/api/", apiRouter)

	fmt.Printf("server started, listening on port %d...\n", conf.ListenPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", conf.ListenPort), nil); err != nil {
		fmt.Printf("failed to start http server on port %d, err: %s\n", conf.ListenPort, err)
		os.Exit(8)
	}
}
