package controllers

import (
	"fmt"
	"github.com/rabobank/scheduler-service-broker/db"
	"github.com/rabobank/scheduler-service-broker/model"
	"github.com/rabobank/scheduler-service-broker/util"
	"net/http"
)

func JobHistoriesGet(w http.ResponseWriter, r *http.Request) {
	if isValid, _, req := ValidateRequest(w, r); isValid {
		if result, err := db.GetJobHistories(req.SpaceGUID, req.Name); err != nil {
			util.WriteHttpResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to list histories: %s", err))
		} else {
			util.WriteHttpResponse(w, http.StatusOK, model.HistoryListResponse{Histories: result})
		}
	}
}
func CallHistoriesGet(w http.ResponseWriter, r *http.Request) {
	if isValid, _, req := ValidateRequest(w, r); isValid {
		if result, err := db.GetCallHistories(req.SpaceGUID, req.Name); err != nil {
			util.WriteHttpResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to list histories: %s", err))
		} else {
			util.WriteHttpResponse(w, http.StatusOK, model.HistoryListResponse{Histories: result})
		}
	}
}
