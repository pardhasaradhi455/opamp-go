package models

type AgentResponse struct {
	ServerName        string `json:"server_name"`
	OS                string `json:"os"`
	Environment       string `json:"environment"`
	OsVersion         string `json:"os_version"`
	AgentVersion      string `json:"agent_version"`
	SupervisorVersion string `json:"supervisor_version"`
	Active            bool   `json:"active"`
}
