package controllers

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/context"
	"github.com/rabobank/scheduler-service-broker/conf"
	"github.com/rabobank/scheduler-service-broker/model"
	"github.com/rabobank/scheduler-service-broker/util"
	"io/ioutil"
	"net/http"
	"time"
)

func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if util.BasicAuth(w, r, conf.BrokerUser, conf.BrokerPassword) {
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
		}
	})
}

func DebugMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.DumpRequest(r)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func AddHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func CheckJWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if accessToken, err := util.GetAccessTokenFromRequest(r); err == nil {
			var token *jwt.Token
			token, err = jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
				} else {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				uaaPubKeyLocation := token.Header["jku"]
				if !util.IsValidJKU(fmt.Sprintf("%s", uaaPubKeyLocation)) {
					return nil, fmt.Errorf("invalid jku parameter in JWT")
				}
				util.PrintfIfDebug("getting uaa pub key from: %s\n", uaaPubKeyLocation)
				transport := http.Transport{IdleConnTimeout: time.Second}
				client := http.Client{Timeout: time.Duration(5) * time.Second, Transport: &transport}
				var resp *http.Response
				if resp, err = client.Get(fmt.Sprintf("%v", uaaPubKeyLocation)); err != nil {
					return nil, err
				} else {
					if resp == nil {
						return nil, fmt.Errorf("empty response from uaa server while getting /token_keys")
					}
					var bodyBytes []byte
					if bodyBytes, err = ioutil.ReadAll(resp.Body); err != nil {
						return nil, err
					} else {
						var uaaPubKeys model.TokenKeys
						if err = json.Unmarshal(bodyBytes, &uaaPubKeys); err != nil {
							return nil, err
						}
						var publicKey []byte
						for _, tokenKey := range uaaPubKeys.Keys {
							if tokenKey.Kid == token.Header["kid"] {
								publicKey = []byte(tokenKey.Value)
							}
						}
						var pubKey *rsa.PublicKey
						if pubKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKey); err != nil {
							return nil, err
						} else {
							return pubKey, nil
						}
					}
				}
			})
			if err != nil {
				fmt.Printf("failed to validate accessToken: %s\n", err)
			} else {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					for mapKey := range claims {
						if mapKey == "user_name" {
							context.Set(r, "jwt", *token) // we use it in subsequent handlers
							util.PrintfIfDebug("successful login for user %s\n", claims[mapKey])
						}
					}
					// Call the next handler, which can be another middleware in the chain, or the final handler.
					next.ServeHTTP(w, r)
					return
				} else {
					fmt.Println("access token is invalid")
				}
			}
		} else {
			fmt.Println("access token is missing")
		}
		w.WriteHeader(401)
		_, _ = w.Write([]byte("Unauthorised.\n"))
	})
}
