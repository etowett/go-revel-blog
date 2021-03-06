package controllers

import "github.com/revel/revel"

type (
	// Response struct for response json
	Response struct {
		App
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
		Status  string      `json:"status"`
	}
)

func response(data interface{}, message string, status string) Response {
	result := Response{}
	result.Message = message
	result.Data = data
	result.Status = status
	return result
}

func init() {
	revel.InterceptMethod(App.AddUser, revel.BEFORE)

	revel.TemplateFuncs["preview"] = func(body string, num int) string {
		return body[0:num]
	}
}
