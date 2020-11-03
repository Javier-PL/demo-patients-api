package models

import "net/http"

type Route struct {
	Path     string                                         `json:"path"`
	Function func(w http.ResponseWriter, req *http.Request) `json:"function"`
	Method   string                                         `json:"method"`
	Mw       string                                         `json:"mw"`
}
