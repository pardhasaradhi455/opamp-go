package response

import (
	"encoding/json"
	"net/http"

	"github.com/open-telemetry/opamp-go/opamp-server/login/logger"
)

type Response struct {
	writer  http.ResponseWriter
	request *http.Request
	status  int
	message string
}

func New(w http.ResponseWriter, r *http.Request) *Response {
	return &Response{
		writer:  w,
		request: r,
	}
}

func (r *Response) Status(status int) *Response {
	r.status = status
	return r
}

func (r *Response) Log(message string) *Response {
	logger.Log(r.status, r.writer, r.request, message)
	r.message = message
	return r
}

func (r *Response) JSON(data interface{}) {
	r.writer.Header().Set("Content-Type", "application/json")
	r.writer.WriteHeader(r.status)

	if err := json.NewEncoder(r.writer).Encode(data); err != nil {
		http.Error(r.writer, err.Error(), http.StatusInternalServerError)
	}
}
