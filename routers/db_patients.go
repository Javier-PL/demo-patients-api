package routers

import (
	"ccl/ccl-patients-api/auth"
	"ccl/ccl-patients-api/models"
	"ccl/ccl-patients-api/services"

	"github.com/gorilla/mux"
)

var routesDBpatients = []models.Route{
	{Path: "/patient/c", Function: services.PostPatient, Method: "POST", Mw: "auth"},
	{Path: "/patient/u", Function: services.UpdatePatient, Method: "POST", Mw: "auth"},
	{Path: "/patient/g", Function: services.GetPatient, Method: "POST", Mw: "auth"},
	{Path: "/patients/g", Function: services.GetPatients, Method: "GET", Mw: "auth"},
	//{Path: "/patient/d", Function: services.DeletePatient, Method: "POST", Mw: "auth"}, //UNUSED
}

func SetRoutesDBpatients(router *mux.Router) *mux.Router {
	for _, r := range routesDBpatients {

		if r.Mw == "" {
			router.HandleFunc(r.Path, r.Function).Methods(r.Method)
		} else if r.Mw == "auth" {
			router.HandleFunc(r.Path, auth.RequireTokenAuthentication(r.Function)).Methods(r.Method)
		}

	}

	return router
}

func GetRoutesDBpatients() []models.Route {
	return routesDBpatients
}
