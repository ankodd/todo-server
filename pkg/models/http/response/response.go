package response

import "net/http"

type Response struct {
	Writer http.ResponseWriter `json:"-"`
	Data   interface{}         `json:"data,omitempty"`
	Status int                 `json:"status"`
	Err    error               `json:"error,omitempty"`
}

func New(object interface{}, status int, err error, w http.ResponseWriter) Response {
	return Response{
		Writer: w,
		Data:   object,
		Status: status,
		Err:    err,
	}
}

func (r *Response) MakeErr(err error, status int) {
	r.Status = status
	r.Err = err

	if r.Data != nil {
		r.Data = nil
	}
}
