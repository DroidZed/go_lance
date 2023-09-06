package db

import (
	"context"
	"time"

	"github.com/DroidZed/go_lance/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbHandle *mongo.Client

func GetConnection() *mongo.Client {

	log := config.InitializeLogger().LogHandler

	if dbHandle != nil {
		return dbHandle
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	dbConnection, err := mongo.Connect(ctx, options.Client().ApplyURI(config.EnvDbURI()))
	if err != nil {
		log.Fatal("Could not connect to database.")
	}

	// Ping the MongoDB server.
	err = dbConnection.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping the database:", err)
	}

	log.Infof("Connected to %s\n", config.EnvDbName())

	dbHandle = dbConnection

	return dbHandle
}
