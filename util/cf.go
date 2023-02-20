package util

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/golang-jwt/jwt"
	"github.com/rabobank/scheduler-service-broker/conf"
	"github.com/rabobank/scheduler-service-broker/model"
	"io/ioutil"
	"os"
	"time"
)

func GetCFClient() *cfclient.Client {
	var err error
	c := &cfclient.Config{
		ApiAddress:        conf.CfApiURL,
		ClientID:          conf.ClientId,
		ClientSecret:      conf.ClientSecret,
		SkipSslValidation: true,
	}
	fmt.Printf("getting cf client from %s...", conf.CfApiURL)
	var client *cfclient.Client
	if client, err = cfclient.NewClient(c); err != nil {
		fmt.Printf("\nfailed getting cf client:%s\n", err)
		os.Exit(8)
	} else {
		fmt.Println("done")
		// refresh the client every hour to get a new refresh token
		go func() {
			channel := time.Tick(time.Duration(conf.TokenRefreshInterval) * time.Minute)
			for range channel {
				if client, err = cfclient.NewClient(c); err != nil {
					fmt.Printf("failed to refresh cf client, error is %s\n", err)
				} else {
					fmt.Println("refreshed cf client, got new token")
					CfClient = *client
				}
			}
		}()
	}
	return client
}

// IsUserAuthorisedForSpace - It takes the jwt, extracts the userId from it,
//
//	then queries cf (/v3/roles) to check if that user has at least developer or manager role for the give space
func IsUserAuthorisedForSpace(token jwt.Token, spaceGuid string) bool {
	userId := token.Claims.(jwt.MapClaims)["user_id"].(string)
	scopes := token.Claims.(jwt.MapClaims)["scope"].([]interface{})
	if Contains(scopes, "cloud_controller.admin") {
		return true
	}
	req := CfClient.NewRequest("GET", fmt.Sprintf("/v3/roles?types=space_developer,space_manager&space_guids=%s&user_guids=%s", spaceGuid, userId))
	if resp, err := CfClient.DoRequest(req); err != nil {
		fmt.Printf("failed to query Cloud Controller for roles: %s\n", err)
		return false
	} else {
		var body []byte
		if body, err = ioutil.ReadAll(resp.Body); err != nil {
			fmt.Printf("failed to read response from /v3/roles query to Cloud Controller: %s\n", err)
			return false
		} else {
			var v3RolesResponse model.GenericV3Response
			if err = json.Unmarshal(body, &v3RolesResponse); err != nil {
				fmt.Printf("failed to parse response from /v3/roles query to Cloud Controller: %s\n", err)
				return false
			} else {
				PrintfIfDebug("found %d roles for userId %s and spaceguid %s\n", v3RolesResponse.Pagination.TotalResults, userId, spaceGuid)
				if v3RolesResponse.Pagination.TotalResults == 0 {
					return false
				}
				return true
			}
		}
	}
}

func IsAppBoundToSchedulerService(appguid string) bool {
	req := CfClient.NewRequest("GET", "/v3/service_plans?service_offering_names=scheduler") // TODO make the service name configurable?
	if resp, err := CfClient.DoRequest(req); err != nil {
		fmt.Println(err)
		return false
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		response := model.GenericV3Response{}
		if err = json.Unmarshal(body, &response); err != nil {
			fmt.Println(err)
			return false
		} else {
			planguid := response.Resources[0].Guid
			req = CfClient.NewRequest("GET", fmt.Sprintf("/v3/service_credential_bindings?app_guids=%s&service_plan_guids=%s", appguid, planguid))
			if resp, err = CfClient.DoRequest(req); err != nil {
				fmt.Println(err)
			} else {
				body, _ = ioutil.ReadAll(resp.Body)
				if err = json.Unmarshal(body, &response); err != nil {
					fmt.Println(err)
					return false
				} else {
					if len(response.Resources) == 1 {
						return true
					}
				}
			}
		}
	}
	return false
}
