package models

type HttpLog struct {
	Path   string
	Method string
	Status int
	Message string
}