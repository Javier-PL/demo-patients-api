package services

import (
	"ccl/ccl-patients-api/database"
	"ccl/ccl-patients-api/models"
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var logs_database = "ccl"
var logs_collection = "ccl.logs"

func PostLog(action string, object string, result string, itemjson string, updatejson string, itemid string) {

	var item models.DBLog
	item.Action = action
	item.Object = object
	item.Result = result
	item.Itemjson = itemjson
	item.Updatejson = updatejson
	item.Date = time.Now()
	item.ItemID = itemid

	m := mymodel{I: item, DatabaseName: logs_database, Collection: logs_collection}

	_, err := createLogInterface(m)
	if err != nil {
		Log.Error(err)
		return
	}

}

func createLogInterface(m mymodel) (*mongo.InsertOneResult, error) {

	c := database.DBCon.Database(m.DatabaseName).Collection(m.Collection)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	res, err := c.InsertOne(ctx, m.I)
	if err != nil {
		Log.Error(err)
		//http.Error(w, err.Error(), 500)
		return nil, err
	}
	//id := res.InsertedID

	return res, nil
}

func GetLogs(w http.ResponseWriter, r *http.Request) {

	var item models.DBLog
	//_ = json.NewDecoder(r.Body).Decode(&item)
	filter := bson.M{}

	itemmodel := mymodel{I: &item, F: filter, DatabaseName: patients_database, Collection: patients_collection}

	itemmodel.getInterfaces(w, r)
}
