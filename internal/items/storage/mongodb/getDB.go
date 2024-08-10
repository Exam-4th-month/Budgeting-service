package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"budgeting-service/internal/items/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(config *config.Config) (*mongo.Database, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		config.MongoDb.User,
		config.MongoDb.Password,
		config.MongoDb.Host,
		config.MongoDb.Port)
	fmt.Println(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(config.MongoDb.DBName)

	log.Printf("--------------------------- Connected to the database %s --------------------------------\n", config.MongoDb.DBName)

	return db, nil
}
