package util

import (
	"fmt"
	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/cloudfoundry/go-cfclient/v3/config"
	"github.com/golang-jwt/jwt"
	"github.com/rabobank/scheduler-service-broker/conf"
	"log"
	"time"
)

func InitCFClient() *client.Client {
	var err error
	if conf.CfConfig, err = config.New(conf.CfApiURL, config.ClientCredentials(conf.ClientId, conf.ClientSecret), config.SkipTLSValidation()); err != nil {
		log.Fatalf("failed to create new config: %s", err)
	}
	if conf.CfClient, err = client.New(conf.CfConfig); err != nil {
		log.Fatalf("failed to create new client: %s", err)
	} else {
		// refresh the client every hour to get a new refresh token
		go func() {
			channel := time.Tick(time.Duration(15) * time.Minute)
			for range channel {
				conf.CfClient, err = client.New(conf.CfConfig)
				if err != nil {
					log.Printf("failed to refresh cfclient, error is %s", err)
				}
			}
		}()
	}
	return conf.CfClient
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
	roleListOptions := client.RoleListOptions{
		ListOptions: &client.ListOptions{},
		Types:       client.Filter{Values: []string{"space_developer", "space_manager"}},
		SpaceGUIDs:  client.Filter{Values: []string{spaceGuid}},
		UserGUIDs:   client.Filter{Values: []string{userId}},
	}
	if roles, err := conf.CfClient.Roles.ListAll(conf.CfCtx, &roleListOptions); err != nil {
		fmt.Printf("failed to query Cloud Controller for roles: %s\n", err)
		return false
	} else {
		PrintfIfDebug("found %d roles for userId %s and spaceguid %s\n", len(roles), userId, spaceGuid)
		if len(roles) == 0 {
			return false
		}
		return true
	}
}

func IsAppBoundToSchedulerService(appguid string) bool {
	planListOptions := client.ServicePlanListOptions{
		ListOptions:          &client.ListOptions{},
		ServiceOfferingNames: client.Filter{Values: []string{"scheduler"}}, // TODO make the service name configurable?
	}
	if plans, err := conf.CfClient.ServicePlans.ListAll(conf.CfCtx, &planListOptions); err != nil {
		fmt.Println(err)
		return false
	} else {
		if len(plans) != 1 {
			return false
		}
		bindingListOptions := client.ServiceCredentialBindingListOptions{
			ListOptions:      &client.ListOptions{},
			AppGUIDs:         client.Filter{Values: []string{appguid}},
			ServicePlanGUIDs: client.Filter{Values: []string{plans[0].GUID}},
		}
		if bindings, err := conf.CfClient.ServiceCredentialBindings.ListAll(conf.CfCtx, &bindingListOptions); err != nil {
			fmt.Println(err)
		} else {
			if len(bindings) == 1 {
				return true
			}
		}
	}
	return false
}
