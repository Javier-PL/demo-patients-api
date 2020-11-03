package services

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"ccl/ccl-patients-api/models"
	"encoding/json"
	"net/http"
)

var patients_database = "ccl"
var patients_collection = "ccl.patients"

func PostPatient(w http.ResponseWriter, r *http.Request) {

	var item models.Patient

	itemmodel := mymodel{I: item, DatabaseName: patients_database, Collection: patients_collection}
	res, err := itemmodel.postInterface(w, r)
	if err != nil {
		Log.Error(err)
		return
	}

	Log.Info("POST", "PATIENT", "ID:", res.InsertedID.(primitive.ObjectID).Hex(), "Title:", "", "SUCCESS:", IsObjectIDValid(res.InsertedID.(primitive.ObjectID)))

}

func GetPatient(w http.ResponseWriter, r *http.Request) {

	var item models.Patient
	_ = json.NewDecoder(r.Body).Decode(&item)

	var filter interface{}

	filter = bson.M{"_id": item.ID}

	itemmodel := mymodel{I: &item, F: filter, DatabaseName: patients_database, Collection: patients_collection}

	itemmodel.getInterface(w, r)
}

func GetPatients(w http.ResponseWriter, r *http.Request) {

	var item models.Patient
	//_ = json.NewDecoder(r.Body).Decode(&item)
	filter := bson.M{}

	itemmodel := mymodel{I: &item, F: filter, DatabaseName: patients_database, Collection: patients_collection}

	itemmodel.getInterfaces(w, r)
}

func UpdatePatient(w http.ResponseWriter, r *http.Request) {

	var filter_update []models.Patient //awaits 2 items, the item to update, and the update
	_ = json.NewDecoder(r.Body).Decode(&filter_update)

	if len(filter_update) < 2 {
		Log.Error(filter_update)
		http.Error(w, "Expecting 2 structures, the original and the update", 500)
		return
	}
	filter := bson.M{"_id": filter_update[0].ID}
	update := bson.M{"$set": filter_update[1]}

	itemmodel := mymodel{I: update, F: filter, DatabaseName: patients_database, Collection: patients_collection}

	res, err := itemmodel.putInterface(w, r)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	Log.Info("UPDATE", "PATIENT", "ID:", filter_update[0].ID, "SUCCESS:", res.ModifiedCount)
	//Log.Info("UPDATE", "INVOICE", "ID:", items[0].ID, "SUCCESS:", IsObjectIDValid(res.UpsertedID.(primitive.ObjectID)))
}

/* //NOT USED
func DeletePatient(w http.ResponseWriter, r *http.Request) {

	var item models.Patient
	itemmodel := mymodel{I: &item, DatabaseName: patients_database, Collection: patients_collection}
	res, err := itemmodel.deleteInterface(w, r)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}
	Log.Info("DELETE", "ID:", item.ID, "Name:", item.Name, "SUCCESS:", res.DeletedCount)
}
*/
