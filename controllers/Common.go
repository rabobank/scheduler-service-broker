package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/context"
	"github.com/rabobank/scheduler-service-broker/util"
	"io/ioutil"
	"net/http"
)

// GenericRequestFitsAll - All request bodies can have different structures, but this one will capture all fields, while some of them may be empty
type GenericRequestFitsAll struct {
	SpaceGUID      string `json:"spaceguid"`
	AppGUID        string `json:"appguid,omitempty"`
	Name           string `json:"name,omitempty"`
	CronExpression string `json:"cronexpression,omitempty"`
	ExpressionType string `json:"expressiontype,omitempty"`
	Command        string `json:"command,omitempty"`
	Url            string `json:"url,omitempty"`
	AuthHeader     string `json:"authheader,omitempty"`
	ScheduleGuid   string `json:"scheduleguid,omitempty"`
}

// ValidateRequest - We validate the incoming http request, it should have a valid JWT, there should be a user_id claim in the JWT, the request body should be json-parse-able and the user should be authorized for the requested space.
func ValidateRequest(w http.ResponseWriter, r *http.Request) (bool, string, GenericRequestFitsAll) {
	var userId string
	var requestObject GenericRequestFitsAll
	if token, ok := context.Get(r, "jwt").(jwt.Token); !ok {
		util.WriteHttpResponse(w, http.StatusBadRequest, "failed to parse access token")
	} else {
		userId = token.Claims.(jwt.MapClaims)["user_id"].(string)
		if body, err := ioutil.ReadAll(r.Body); err != nil {
			util.WriteHttpResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to read request body: %s", err))
		} else {
			if err = json.Unmarshal(body, &requestObject); err != nil {
				util.WriteHttpResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to parse request body: %s", err))
			} else {
				if util.IsUserAuthorisedForSpace(token, requestObject.SpaceGUID) {
					return true, userId, requestObject
				} else {
					util.WriteHttpResponse(w, http.StatusUnauthorized, fmt.Sprintf("you are not authorized for space %s", requestObject.SpaceGUID))
				}
			}
		}
	}
	return false, userId, requestObject
}
