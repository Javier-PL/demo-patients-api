package routers

import (
	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router = SetRoutesDBinvoices(router)
	router = SetRoutesDBpatients(router)
	router = SetRoutesDBlogs(router)

	return router
}
