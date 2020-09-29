package models

// RequestPublishNews struct is request models for publish news api
type RequestPublishNews struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}

// ResponseNews struct is response models for  news api
type ResponseNews struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
