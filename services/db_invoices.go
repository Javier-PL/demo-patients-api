package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"ccl/ccl-patients-api/models"
	"encoding/json"
	"net/http"
	"sync"
)

var invoices_database = "ccl"
var invoices_collection = "ccl.invoices"

type CounterInvoices struct {
	sync.Mutex
	Total int64
}

var TotalDBInvoices CounterInvoices

func PostInvoice(w http.ResponseWriter, r *http.Request) {

	var item models.Invoice

	itemmodel := mymodel{I: item, DatabaseName: invoices_database, Collection: invoices_collection}
	res, err := itemmodel.postInterface(w, r)
	if err != nil {
		Log.Error(err)
		return
	}
	Log.Info("POST", "INVOICE", "ID:", res.InsertedID.(primitive.ObjectID).Hex(), "SUCCESS:", IsObjectIDValid(res.InsertedID.(primitive.ObjectID)))
}

func (counter *CounterInvoices) PostSyncInvoice(w http.ResponseWriter, r *http.Request) {

	counter.Lock()
	defer counter.Unlock()

	var item models.Invoice
	_ = json.NewDecoder(r.Body).Decode(&item)

	IN, _ := strconv.ParseInt(os.Getenv("INVOICE_NUMBER"), 10, 64)

	item.InvoiceNumber = counter.Total + IN

	itemmodel := mymodel{I: item, DatabaseName: invoices_database, Collection: invoices_collection}

	res, err := createInterface(itemmodel)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	respBody, err := json.Marshal(res)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	counter.Total++

	Log.Info("POST SYNC", "INVOICE", "ID:", res.InsertedID.(primitive.ObjectID).Hex(), "SUCCESS:", IsObjectIDValid(res.InsertedID.(primitive.ObjectID)), "with counter:", TotalDBInvoices.Total)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(respBody)

}

func (counter *CounterInvoices) GetInvoiceNumber(w http.ResponseWriter, r *http.Request) {

	counter.Lock()
	defer counter.Unlock()

	IN, _ := strconv.ParseInt(os.Getenv("INVOICE_NUMBER"), 10, 64)

	InvoiceNumber := counter.Total + IN

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(`{"InvoiceNumber":"` + strconv.FormatInt(InvoiceNumber, 10) + `"}`))

}

/* //UNUSED
func GetInvoice(w http.ResponseWriter, r *http.Request) {

	var item models.Invoice
	_ = json.NewDecoder(r.Body).Decode(&item)

	var filter interface{}

	filter = bson.M{"_id": item.ID}

	itemmodel := mymodel{I: &item, F: filter, DatabaseName: invoices_database, Collection: invoices_collection}

	itemmodel.getInterface(w, r)
}
*/

func GetInvoices(w http.ResponseWriter, r *http.Request) {

	keys, _ := r.URL.Query()["patientID"]
	key := keys[0]

	var item models.Invoice

	var filter interface{}
	if key != "" {
		filter = bson.M{"patientid": key}
	} else {
		filter = bson.M{}

	}

	itemmodel := mymodel{I: &item, F: filter, DatabaseName: invoices_database, Collection: invoices_collection}

	itemmodel.getInterfaces(w, r)
}

func GetInvoiceByInvoiceNumber(w http.ResponseWriter, r *http.Request) {

	keys, _ := r.URL.Query()["invoice_number"]
	key := keys[0]

	var item models.Invoice

	var filter interface{}
	if key != "" {
		key64, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			Log.Error(err)
			http.Error(w, err.Error(), 500)
			return
		}
		filter = bson.M{"invoice_number": key64}
	} else {
		filter = bson.M{}

	}
	Log.Info(filter)
	itemmodel := mymodel{I: &item, F: filter, DatabaseName: invoices_database, Collection: invoices_collection}

	itemmodel.getInterface(w, r)
}

func GetInvoicesRange(w http.ResponseWriter, r *http.Request) {

	var item models.Invoice

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	layout := "2006-01-02T15:04:05.000 00:00"
	fromDate, err := time.Parse(layout, from)
	if err != nil {
		fmt.Println("error", err)
	}

	toDate, err := time.Parse(layout, to)
	if err != nil {
		fmt.Println(err)
	}

	fromDate = fromDate.AddDate(0, 0, -1) //the selected dates come from frontend, but when filtering them, golang doesnt include these two dates in the result. So take 1 day before and 1 after.
	toDate = toDate.AddDate(0, 0, 1)

	filter := bson.M{"date": bson.M{
		"$gt": fromDate, //time.Date(2020, time.April, 20, 0, 0, 0, 0, time.UTC),
		"$lt": toDate,   //time.Date(2020, time.April, 22, 0, 0, 0, 0, time.UTC),
	}}

	fmt.Println(fromDate, toDate)

	itemmodel := mymodel{I: &item, F: filter, DatabaseName: invoices_database, Collection: invoices_collection}

	itemmodel.getInterfaces(w, r)
}

/*
func UpdateInvoice(w http.ResponseWriter, r *http.Request) {

	var filter_update []models.Invoice //awaits 2 items, the item to update, and the update
	_ = json.NewDecoder(r.Body).Decode(&filter_update)

	if len(filter_update) < 2 {
		Log.Error(filter_update)
		http.Error(w, "Expecting 2 structures, the original and the update", 500)
		return
	}

	filter := bson.M{"_id": filter_update[0].ID}
	update := bson.M{"$set": filter_update[1]}

	itemmodel := mymodel{I: update, F: filter, DatabaseName: invoices_database, Collection: invoices_collection}

	res, err := itemmodel.putInterface(w, r)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	Log.Info("UPDATE", "INVOICE", "ID:", filter_update[0].ID, "SUCCESS:", res.ModifiedCount)
}*/

func UpdateInvoice(w http.ResponseWriter, r *http.Request) {

	//filter := bson.M{"_id": filter_update[0].ID}
	//update := bson.M{"$set": filter_update[1]}

	keys, _ := r.URL.Query()["ID"]
	key := keys[0]

	var updateData models.Invoice
	_ = json.NewDecoder(r.Body).Decode(&updateData)
	update := bson.M{"$set": updateData}

	var filter interface{}
	if key != "" {
		idPrimitive, err := primitive.ObjectIDFromHex(key)
		if err != nil {
			Log.Error(err)
			http.Error(w, err.Error(), 500)
			return
		}
		filter = bson.M{"_id": idPrimitive}
	} else {
		filter = bson.M{}

	}

	itemmodel := mymodel{I: update, F: filter, DatabaseName: invoices_database, Collection: invoices_collection}

	res, err := itemmodel.putInterface(w, r)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	Log.Info("UPDATE", "INVOICE", "ID:", filter, "SUCCESS:", res.ModifiedCount)
}

func UpdateInvoicesUserData(w http.ResponseWriter, r *http.Request) {

	var filter_update []models.Invoice //awaits 2 items, the item to update, and the update
	_ = json.NewDecoder(r.Body).Decode(&filter_update)

	if len(filter_update) < 2 {
		Log.Error(filter_update)
		http.Error(w, "Expecting 2 structures, the filter and the update", 500)
		return
	}
	//filter := bson.M{"patientid": items[0].PatientID}

	var filter interface{}

	if filter_update[0].PatientID != "" {
		filter = bson.M{"patientid": filter_update[0].PatientID}
	} else {
		filter = bson.M{"_id": filter_update[0].ID}
	}

	var update interface{}
	var partialUpdate interface{}

	partialUpdate = bson.M{
		"patient":        filter_update[1].Patient,
		"patientdni":     filter_update[1].PatientDNI,
		"patientaddress": filter_update[1].PatientAddress}

	update = bson.M{"$set": partialUpdate}

	itemmodel := mymodel{I: update, F: filter, DatabaseName: invoices_database, Collection: invoices_collection}

	res, err := itemmodel.putManyInterface(w, r)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	Log.Info("UPDATE", "INVOICES", "ID:", filter_update[0].PatientID, "SUCCESS:", res.ModifiedCount)
}

func (counter *CounterInvoices) CountInvoices() error {

	itemmodel := mymodel{I: nil, F: nil, DatabaseName: invoices_database, Collection: invoices_collection}

	result, err := countInterfaces(itemmodel)
	if err != nil {
		Log.Error(err)
		return err
	}

	counter.Total = result

	return nil

}

/* //NOT USED
func DeleteInvoice(w http.ResponseWriter, r *http.Request) {

	var item models.Invoice
	itemmodel := mymodel{I: &item, DatabaseName: invoices_database, Collection: invoices_collection}
	res, err := itemmodel.deleteInterface(w, r)
	if err != nil {
		Log.Error(err)
		return
	}
	Log.Info("DELETE", "INVOICE", "ID:", item.ID, "SUCCESS:", res.DeletedCount)
}
*/
