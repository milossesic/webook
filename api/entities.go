package api

// swagger:model myfirstapi
//type notification
//struct {
//	Msg string
//}

// Returns internal error when something went wrong on the server side
// swagger:model internalError
type internalError struct {
	// example: 500
	ErrorCode int `json:"code"`
	// default: Uh oh! Something went wrong on the server side
	Message string `json:"msg"`
	// default: reason for failure
	Cause string `json:"cause"`
}
