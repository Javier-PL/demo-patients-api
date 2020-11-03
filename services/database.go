package services

import (
	"reflect"
	"strconv"

	"github.com/Ramso-dev/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"ccl/ccl-patients-api/database"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

var Log log.Logger

type mymodel struct {
	I            interface{}
	F            interface{}
	DatabaseName string
	Collection   string
}

func (m mymodel) postInterface(w http.ResponseWriter, r *http.Request) (*mongo.InsertOneResult, error) {

	//var attendace models.UserAttendance
	_ = json.NewDecoder(r.Body).Decode(&m.I)

	res, err := createInterface(m)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return nil, err
	}

	respBody, err := json.Marshal(res)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return nil, err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(respBody)

	return res, nil
}

func createInterface(m mymodel) (*mongo.InsertOneResult, error) {

	c := database.DBCon.Database(m.DatabaseName).Collection(m.Collection)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	tolog, err := json.Marshal(m.I)
	if err != nil {
		Log.Error(err)
	}

	res, err := c.InsertOne(ctx, m.I)
	if err != nil {
		Log.Error(err)
		go PostLog("POST", m.Collection, "ERROR", string(tolog), "", res.InsertedID.(primitive.ObjectID).String())
		return nil, err
	}
	//id := res.InsertedID

	go PostLog("POST", m.Collection, "SUCCESS", string(tolog), "", res.InsertedID.(primitive.ObjectID).String())

	return res, nil
}

func (m mymodel) getInterface(w http.ResponseWriter, r *http.Request) {

	result, err := readInterface(m)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	respBody, err := json.Marshal(&result)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(respBody)

}

func readInterface(m mymodel) (*interface{}, error) {

	c := database.DBCon.Database(m.DatabaseName).Collection(m.Collection)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err := c.FindOne(ctx, m.F).Decode(m.I)
	if err != nil {
		Log.Error(err)
		return nil, err
	}

	return &m.I, nil

}

func (m mymodel) getInterfaces(w http.ResponseWriter, r *http.Request) {

	result, err := readInterfaces(m)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	respBody, err := json.Marshal(result)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	//Log.Info(string(respBody))

	//time.Sleep(4 * time.Second)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(respBody)

}

func readInterfaces(m mymodel) (*[]interface{}, error) {

	c := database.DBCon.Database(m.DatabaseName).Collection(m.Collection)

	var results []interface{}

	// Pass these options to the Find method
	options := options.Find()
	options.SetLimit(2000)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := c.Find(ctx, m.F, options)

	if err != nil {
		Log.Error(err)
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {

		//create a new struct from the struct pointer
		rv := reflect.ValueOf(&m.I).Elem()
		typ := rv.Elem().Type().Elem()
		rv.Set(reflect.New(typ))
		rStruct := rv.Interface()

		err := cur.Decode(rStruct)
		if err != nil {
			Log.Error(err)
			return nil, err
		}

		results = append(results, rStruct)

	}
	if err := cur.Err(); err != nil {
		Log.Error(err)
		return nil, err
	}

	return &results, nil

}

func (m mymodel) putInterface(w http.ResponseWriter, r *http.Request) (*mongo.UpdateResult, error) {

	result, err := updateInterface(m)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return nil, err
	}

	respBody, err := json.Marshal(result)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return nil, err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(respBody)

	return result, nil
}
func updateInterface(m mymodel) (*mongo.UpdateResult, error) {

	c := database.DBCon.Database(m.DatabaseName).Collection(m.Collection)

	//update := bson.M{"$set": m.I}

	updatetolog, err := json.Marshal(m.I)
	if err != nil {
		Log.Error(err)
	}
	originaltolog, err := json.Marshal(m.F)
	if err != nil {
		Log.Error(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	options := options.Update()
	options.SetUpsert(false) //SetUpsert(true) //create if not existing

	updateResult, err := c.UpdateOne(ctx, m.F, m.I, options)
	if err != nil {
		Log.Error(err)
		go PostLog("PUT", m.Collection, "ERROR", string(originaltolog), string(updatetolog), "")
		return nil, err
	}

	go PostLog("PUT", m.Collection, "SUCCESS 1", string(originaltolog), string(updatetolog), "")

	return updateResult, nil

}

func (m mymodel) putManyInterface(w http.ResponseWriter, r *http.Request) (*mongo.UpdateResult, error) {

	result, err := updateManyInterface(m)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return nil, err
	}

	respBody, err := json.Marshal(result)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return nil, err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(respBody)

	return result, nil
}

func updateManyInterface(m mymodel) (*mongo.UpdateResult, error) {

	c := database.DBCon.Database(m.DatabaseName).Collection(m.Collection)

	updatetolog, err := json.Marshal(m.I)
	if err != nil {
		Log.Error(err)
	}
	originaltolog, err := json.Marshal(m.F)
	if err != nil {
		Log.Error(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	options := options.Update()
	options.SetUpsert(false) //SetUpsert(true) //create if not existing

	updateResult, err := c.UpdateMany(ctx, m.F, m.I, options)
	if err != nil {
		Log.Error(err)
		go PostLog("PUT", m.Collection, "ERROR", string(originaltolog), string(updatetolog), "")
		return nil, err
	}

	go PostLog("PUT", m.Collection, "SUCCESS "+strconv.FormatInt(updateResult.ModifiedCount, 10), string(originaltolog), string(updatetolog), "")

	return updateResult, nil

}

func (m mymodel) deleteInterface(w http.ResponseWriter, r *http.Request) (*mongo.DeleteResult, error) {

	_ = json.NewDecoder(r.Body).Decode(&m.I)

	res, err := removeInterface(m)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return nil, err
	}

	respBody, err := json.Marshal(res)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return nil, err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(respBody)

	return res, nil

}

func removeInterface(m mymodel) (*mongo.DeleteResult, error) {

	c := database.DBCon.Database(m.DatabaseName).Collection(m.Collection)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	res, err := c.DeleteOne(ctx, m.I)
	if err != nil {
		Log.Error(err)
		//http.Error(w, err.Error(), 500)
		return nil, err
	}
	//id := res.InsertedID

	return res, nil

}

func countInterfaces(m mymodel) (int64, error) {

	c := database.DBCon.Database(m.DatabaseName).Collection(m.Collection)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	itemCount, err := c.CountDocuments(ctx, bson.M{})
	if err != nil {
		Log.Info(err)
		return -1, err
	}

	return itemCount, nil

}

func (m mymodel) getCountInterfaces(w http.ResponseWriter, r *http.Request) {

	result, err := countInterfaces(m)
	if err != nil {
		Log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	respBody, err := json.Marshal(&result)
	if err != nil {
		Log.Info(err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(respBody)

}

//IsObjectIDEmpty checks if a an ObjectID is valid for our database, cause if it is not send, it is a row of 0s per default
func IsObjectIDValid(objID primitive.ObjectID) bool {
	if objID.Hex() == "000000000000000000000000" {
		return false
	}
	return true
}
