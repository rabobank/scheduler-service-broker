package model

type CallListResponse struct {
	Calls []Call
}

type Call struct {
	CallName   string `json:"callname"`
	AppName    string `json:"appname"`
	Url        string `json:"url"`
	AuthHeader string `json:"authheader"`
}
