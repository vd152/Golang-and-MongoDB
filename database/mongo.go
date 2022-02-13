package database

import (
	"context"
	"log"
	"os"

	env "github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func ConnectDB() {
	err := env.Load()

	if err != nil {
		panic("Cannot read dotenv file.")
	}

	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		panic(err)
	} else {
		log.Println("MongoDB connected")
	}

}

func GetMongoClient() *mongo.Client {
	return client
}
