package routers

import (
	"ccl/ccl-patients-api/auth"
	"ccl/ccl-patients-api/models"
	"ccl/ccl-patients-api/services"

	"github.com/gorilla/mux"
)

var routesDBlogs = []models.Route{

	{Path: "/logs/g", Function: services.GetLogs, Method: "GET", Mw: ""},
}

func SetRoutesDBlogs(router *mux.Router) *mux.Router {
	for _, r := range routesDBlogs {

		if r.Mw == "" {
			router.HandleFunc(r.Path, r.Function).Methods(r.Method)
		} else if r.Mw == "auth" {
			router.HandleFunc(r.Path, auth.RequireTokenAuthentication(r.Function)).Methods(r.Method)
		}

	}

	return router
}

func GetRoutesDBlogs() []models.Route {
	return routesDBlogs
}
