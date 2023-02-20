package model

type TokenKeys struct {
	Keys []TokenKey `json:"keys"`
}

type TokenKey struct {
	Kty   string `json:"kty"`
	E     string `json:"e"`
	Use   string `json:"use"`
	Kid   string `json:"kid"`
	Alg   string `json:"alg"`
	Value string `json:"value"`
	N     string `json:"n"`
}
