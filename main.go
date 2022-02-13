package main

import (
	Controller "GoREST/controllers"
	Database "GoREST/database"
	"context"
	"log"
	"net/http"
)

func main() {
	Database.ConnectDB()
	client := Database.GetMongoClient()
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		} else {
			log.Println("DB disconnected.")
		}
	}()
	port := ":3000"

	//CREATE
	http.HandleFunc("/add", Controller.AddRecord)

	//READ
	http.HandleFunc("/", Controller.GetAllRecords)
	http.HandleFunc("/get", Controller.GetRecordById)

	//UPDATE
	http.HandleFunc("/update", Controller.UpdateRecord)

	//DELETE
	http.HandleFunc("/delete", Controller.DeleteRecord)

	log.Println("Server running on", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
