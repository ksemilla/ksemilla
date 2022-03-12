package database

import (
	"context"
	"log"
	"time"

	"ksemilla/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"path/filepath"

	"github.com/joho/godotenv"
)

type DB struct {
	client *mongo.Client
}

func Connect() *DB {

	godotenv.Load(filepath.Join(".", ".env"))
	config := config.Config()

	client, err := mongo.NewClient(options.Client().ApplyURI(config.MONGODB_URI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return &DB{
		client: client,
	}
}
