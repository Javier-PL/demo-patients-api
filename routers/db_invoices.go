package routers

import (
	"ccl/ccl-patients-api/auth"
	"ccl/ccl-patients-api/models"
	"ccl/ccl-patients-api/services"

	"github.com/gorilla/mux"
)

var routesDBinvoices = []models.Route{

	{Path: "/invoice/number", Function: services.TotalDBInvoices.GetInvoiceNumber, Method: "GET", Mw: ""},
	{Path: "/invoice/c", Function: services.TotalDBInvoices.PostSyncInvoice, Method: "POST", Mw: "auth"},
	{Path: "/invoice/u", Function: services.UpdateInvoice, Method: "PUT", Mw: "auth"},
	{Path: "/invoices/u", Function: services.UpdateInvoicesUserData, Method: "POST", Mw: "auth"},
	//{Path: "/invoice/g", Function: services.GetInvoice, Method: "POST", Mw: "auth"}, //UNUSED
	{Path: "/invoice/g/invoice_number", Function: services.GetInvoiceByInvoiceNumber, Method: "GET", Mw: "auth"}, //UNUSED
	{Path: "/invoices/g/range", Function: services.GetInvoicesRange, Method: "GET", Mw: "rangedfilter"},
	{Path: "/invoices/g", Function: services.GetInvoices, Method: "GET", Mw: "auth"},
	//{Path: "/invoice/d", Function: services.DeleteInvoice, Method: "POST", Mw: "auth"}, //UNUSED
	{Path: "/invoices/historial_pdf", Function: services.CreatePDF, Method: "POST", Mw: "auth"},
	{Path: "/invoices/invoice_pdf", Function: services.CreatePatientInvoicePDF, Method: "POST", Mw: "auth"},
}

func SetRoutesDBinvoices(router *mux.Router) *mux.Router {
	for _, r := range routesDBinvoices {

		if r.Mw == "" {
			router.HandleFunc(r.Path, r.Function).Methods(r.Method)

		} else if r.Mw == "rangedfilter" {
			//router.HandleFunc(r.Path, r.Function).Methods(r.Method).Queries("from", "{from}").Queries("to", "{to}")

			router.HandleFunc(r.Path, auth.RequireTokenAuthentication(r.Function)).Methods(r.Method).Queries("from", "{from}").Queries("to", "{to}")

		} else if r.Mw == "auth" {
			router.HandleFunc(r.Path, auth.RequireTokenAuthentication(r.Function)).Methods(r.Method)
		}

	}

	return router
}

func GetRoutesDBinvoices() []models.Route {
	return routesDBinvoices
}
