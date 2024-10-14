package controllers

import (
	"fmt"
	"github.com/rabobank/scheduler-service-broker/conf"
	"github.com/rabobank/scheduler-service-broker/cron"
	"github.com/rabobank/scheduler-service-broker/db"
	"github.com/rabobank/scheduler-service-broker/model"
	"github.com/rabobank/scheduler-service-broker/util"
	"net/http"
	"time"
)

func CallCreate(w http.ResponseWriter, r *http.Request) {
	if isValid, userId, req := ValidateRequest(w, r); isValid {
		if req.Name == "" {
			util.WriteHttpResponse(w, http.StatusBadRequest, "the name of the call was empty")
		} else {
			if _, err := conf.CfClient.Applications.Get(conf.CfCtx, req.AppGUID); err != nil {
				util.WriteHttpResponse(w, http.StatusBadRequest, fmt.Sprintf("app with guid %s not found: %s", req.AppGUID, err))
			} else {
				if !util.IsAppBoundToSchedulerService(req.AppGUID) {
					util.WriteHttpResponse(w, http.StatusBadRequest, fmt.Sprintf("app with guid %s is not bound to an instance of scheduler", req.AppGUID))
				} else {
					// try to insert the call
					if callguid, err := db.InsertCall(db.Call{AppGuid: req.AppGUID, SpaceGuid: req.SpaceGUID, Name: req.Name, Url: req.Url, AuthHeader: req.AuthHeader}); err != nil {
						util.WriteHttpResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to create call: %s", err))
					} else {
						fmt.Printf("userId %s created call with guid %s for space guid %s\n", userId, callguid, req.SpaceGUID)
						util.WriteHttpResponse(w, http.StatusCreated, fmt.Sprintf("call created"))
					}
				}
			}
		}
	}
}

func CallRun(w http.ResponseWriter, r *http.Request) {
	if isValid, userId, req := ValidateRequest(w, r); isValid {
		if req.Name == "" {
			util.WriteHttpResponse(w, http.StatusBadRequest, "the name of the call was empty")
		} else {
			if existingCalls, err := db.GetCalls(req.SpaceGUID, req.Name); err == nil && len(existingCalls) > 0 {
				// there should only be a result of length 1:
				call := existingCalls[0]
				go func() {
					cron.DoCall(time.Now(), model.SchedulableCall{CallName: call.Name, AppGuid: call.AppGuid, SpaceGuid: call.SpaceGuid, Url: call.Url, AuthHeader: call.AuthHeader})
				}()
				fmt.Printf("userId %s ran call for space guid %s\n", userId, req.SpaceGUID)
				util.WriteHttpResponse(w, http.StatusOK, fmt.Sprintf("run scheduled for call %s", req.Name))
			} else {
				util.WriteHttpResponse(w, http.StatusBadRequest, fmt.Sprintf("call %s was not found", req.Name))
			}
		}
	}
}

func CallGet(w http.ResponseWriter, r *http.Request) {
	if isValid, _, req := ValidateRequest(w, r); isValid {
		if result, err := db.GetCalls(req.SpaceGUID, ""); err != nil {
			util.WriteHttpResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to list calls: %s", err))
		} else {
			var calls = make([]model.Call, 0)
			for _, call := range result {
				appName := "<unknown>"
				if app, err := conf.CfClient.Applications.Get(conf.CfCtx, call.AppGuid); err == nil {
					appName = app.Name
				}
				calls = append(calls, model.Call{CallName: call.Name, AppName: appName, Url: call.Url, AuthHeader: call.AuthHeader})
			}
			util.WriteHttpResponse(w, http.StatusOK, model.CallListResponse{Calls: calls})
		}
	}
}

func CallDelete(w http.ResponseWriter, r *http.Request) {
	var err error
	if isValid, userId, req := ValidateRequest(w, r); isValid {
		if req.Name == "" {
			util.WriteHttpResponse(w, http.StatusBadRequest, "the name of the call was empty")
		} else {
			fmt.Printf("deleting call with name %s for user %s for spaceguid %s from client %s...\n", req.Name, userId, req.SpaceGUID, r.RemoteAddr)
			if err = db.DeleteCallBySpaceGuidAndCallname(req.SpaceGUID, req.Name); err != nil {
				util.WriteHttpResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete call with name %s in spaceguid %s: %s", req.Name, req.SpaceGUID, err))
			} else {
				cron.DeleteCallByNameAndSpaceGuid(req.Name, req.SpaceGUID)
				util.WriteHttpResponse(w, http.StatusOK, fmt.Sprintf("call with name %s deleted", req.Name))
			}
		}
	}
}
