package util

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rabobank/scheduler-service-broker/conf"
	"github.com/rabobank/scheduler-service-broker/model"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func WriteHttpResponse(w http.ResponseWriter, code int, object interface{}) {
	if data, err := json.Marshal(object); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	} else {
		w.WriteHeader(code)
		_, _ = fmt.Fprintf(w, string(data))
		PrintfIfDebug("response: code:%d, body: %s\n", code, string(data))
	}
}

// BasicAuth validate if user/pass in the http request match the configured service broker user/pass
func BasicAuth(w http.ResponseWriter, r *http.Request, username, password string) bool {
	if user, pass, ok := r.BasicAuth(); !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
		w.Header().Set("WWW-Authenticate", `Basic realm="`+conf.BasicAuthRealm+`"`)
		w.WriteHeader(401)
		_, _ = w.Write([]byte("Unauthorised.\n"))
		return false
	}
	return true
}

func DumpRequest(r *http.Request) {
	if conf.Debug {
		fmt.Printf("dumping %s request for URL: %s\n", r.Method, r.URL)
		fmt.Println("dumping request headers...")
		// Loop over header names
		for name, values := range r.Header {
			if name == "Authorization" {
				fmt.Printf(" %s: %s\n", name, "<redacted>")
			} else {
				// Loop over all values for the name.
				for _, value := range values {
					fmt.Printf(" %s: %s\n", name, value)
				}
			}
		}

		// dump the request body
		fmt.Println("dumping request body...")
		if body, err := io.ReadAll(r.Body); err != nil {
			fmt.Printf("Error reading body: %v\n", err)
		} else {
			fmt.Println(string(body))
			// Restore the io.ReadCloser to it's original state
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}
	}
}

func ProvisionObjectFromRequest(r *http.Request, object interface{}) error {
	if body, err := io.ReadAll(r.Body); err != nil {
		fmt.Printf("failed to read json object from request, error: %s\n", err)
		return err
	} else {
		fmt.Printf("received body:%v\n", string(body))
		if err = json.Unmarshal(body, object); err != nil {
			fmt.Printf("failed to parse json object from request, error: %s\n", err)
			return err
		}
		return nil
	}
}

func GetServiceById(serviceId string) model.Service {
	var service model.Service
	for _, service = range conf.Catalog.Services {
		if service.Id == serviceId {
			return service
		}
	}
	return service
}

// GetAccessTokenFromRequest - get the JWT from the request
func GetAccessTokenFromRequest(r *http.Request) (string, error) {
	var accessToken string
	if authHeaders := r.Header["Authorization"]; authHeaders != nil && len(authHeaders) != 0 {
		accessToken = strings.TrimPrefix(authHeaders[0], "bearer ")
	} else {
		return accessToken, errors.New("no Authorization header found")
	}
	return accessToken, nil
}

func GenerateGUID() string {
	ba := make([]byte, 16)
	if _, err := rand.Read(ba); err != nil {
		fmt.Printf("failed to generate guid, err: %s\n", err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", ba[0:4], ba[4:6], ba[6:8], ba[8:10], ba[10:])
}

func LastXChars(victim string, maxLen int) string {
	if len(victim) > maxLen {
		return victim[len(victim)-maxLen:]
	}
	return victim
}

func PrintIfDebug(msg string) {
	if conf.Debug {
		fmt.Print(msg)
	}
}

func PrintfIfDebug(msg string, args ...interface{}) {
	PrintIfDebug(fmt.Sprintf(msg, args...))
}

// IsValidJKU - We compare the jku with the api hostname, only the first part should be different: like uaa.sys.cfd04.aws.rabo.cloud versus api.sys.cfd04.aws.rabo.cloud
func IsValidJKU(jkuURL string) bool {
	parsedJkuURL, err := url.Parse(jkuURL)
	if err != nil {
		fmt.Printf("jku URL %s is invalid: %s", jkuURL, err)
		return false
	}
	apiURL, _ := url.Parse(conf.CfApiURL)
	apiDomain := strings.TrimPrefix(apiURL.Hostname(), "api.")
	jkuDomain := strings.TrimPrefix(parsedJkuURL.Hostname(), "uaa.")
	if jkuDomain != apiDomain {
		fmt.Printf("jku URL %s is invalid", jkuURL)
		return false
	}
	return true
}

func Contains(elems []interface{}, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
