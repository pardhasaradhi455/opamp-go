package models

type ErrorResponse struct {
	Message string `json:"message"`
}

type TokenErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
