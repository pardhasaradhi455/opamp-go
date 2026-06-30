package logger

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/open-telemetry/opamp-go/opamp-server/login/models"
)

func Log(status int, w http.ResponseWriter, r *http.Request, message string) {
	logEntry := models.HttpLog{Path: r.RequestURI, Method: r.Method, Status: status, Message: message}

	data, _ := json.Marshal(logEntry)
	log.Println(string(data))
}
