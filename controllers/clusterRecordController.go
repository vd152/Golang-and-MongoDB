package controllers

import (
	Database "GoREST/database"
	Model "GoREST/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllRecords(w http.ResponseWriter, r *http.Request) {
	client := Database.GetMongoClient()
	coll := client.Database("myFirstDatabase").Collection("clusters")
	cursor, err := coll.Find(context.TODO(), bson.D{})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
		panic(err)
	}

	var result []bson.M
	if err = cursor.All(context.TODO(), &result); err != nil {
		panic(err)
	}

	for _, result := range result {
		output, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "%s \n", output)
	}
}

func GetRecordById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		fmt.Fprint(w, "ID is required")
		return
	}

	client := Database.GetMongoClient()
	collection := client.Database("myFirstDatabase").Collection("clusters")

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Fprint(w, "Inavlid id")
		return
	}
	var result bson.M
	err = collection.FindOne(context.TODO(), bson.M{"_id": idPrimitive}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Fprint(w, "No documents found.")
			return
		}
	}

	output, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "%s \n", output)
}

func AddRecord(w http.ResponseWriter, r *http.Request) {
	//?name=cluster45&hostname=145.25.65.2&environment=dev
	name, hostname, environment := r.URL.Query().Get("name"), r.URL.Query().Get("hostname"), r.URL.Query().Get("environment")

	if name == "" || hostname == "" || environment == "" {
		fmt.Fprint(w, "Incomplete data.")
		return
	}
	client := Database.GetMongoClient()
	coll := client.Database("myFirstDatabase").Collection("clusters")

	doc := Model.ClusterRecord{
		Name:        name,
		Hostname:    hostname,
		Environment: environment,
		Status:      "created",
	}
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "Document added : %s \n", result.InsertedID)
}

func UpdateRecord(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		fmt.Fprint(w, "ID is required")
		return
	}
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Fprint(w, "Inavlid id")
		return
	}
	name, hostname, environment := r.URL.Query().Get("name"), r.URL.Query().Get("hostname"), r.URL.Query().Get("environment")

	if name == "" && hostname == "" && environment == "" {
		fmt.Fprint(w, "Please add daata to update.")
		return
	}
	client := Database.GetMongoClient()
	collection := client.Database("myFirstDatabase").Collection("clusters")

	var foundDoc = Model.ClusterRecord{}
	err = collection.FindOne(context.TODO(), bson.M{"_id": idPrimitive}).Decode(&foundDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Fprint(w, "No documents found.")
			return
		}
	}

	if name == "" {
		name = foundDoc.Name
	}
	if hostname == "" {
		hostname = foundDoc.Hostname
	}

	if environment == "" {
		environment = foundDoc.Environment
	}
	doc := Model.ClusterRecord{
		Name:        name,
		Hostname:    hostname,
		Environment: environment,
		Status:      foundDoc.Status,
	}

	result, err := collection.ReplaceOne(context.TODO(), bson.M{"_id": idPrimitive}, doc)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprintf(w, "Documents updated: %v\n", result.ModifiedCount)
}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		fmt.Fprint(w, "ID is required")
		return
	}
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Fprint(w, "Inavlid id")
		return
	}
	client := Database.GetMongoClient()
	collection := client.Database("myFirstDatabase").Collection("clusters")
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": idPrimitive})
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "deleted %v documents\n", result.DeletedCount)

}
