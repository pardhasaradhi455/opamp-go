package models

import "net/http"

type RouteDef struct {
	Path            string
	Method          string
	HandlerFunction http.HandlerFunc
	Protected       bool
}
