package response

import "net/http"

// Response is a http response
type Response struct {
	Writer http.ResponseWriter `json:"-"`
	Data   interface{}         `json:"data,omitempty"`
	Status int                 `json:"status"`
	Err    error               `json:"error,omitempty"`
}

// LogResponse is used to log the response
type LogResponse struct {
	Status int
	Err    error
}

// New returns a new Response
func New(object interface{}, status int, err error, w http.ResponseWriter) Response {
	return Response{
		Writer: w,
		Data:   object,
		Status: status,
		Err:    err,
	}
}

// Log returns a LogResponse
func (r *Response) Log() LogResponse {
	return LogResponse{
		Status: r.Status,
		Err:    r.Err,
	}
}

// MakeErr makes an error
func (r *Response) MakeErr(err error, status int) {
	r.Status = status
	r.Err = err

	if r.Data != nil {
		r.Data = nil
	}
}
