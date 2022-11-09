package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func ConnectDatabaseError(conf *MongoConfig) (error, *mongo.Database, *mongo.Client) {

	clientMongo, err := mongo.NewClient(options.Client().ApplyURI(conf.Uri))
	if err != nil {
		return err, nil, nil
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = clientMongo.Connect(ctx)

	if err != nil {
		return err, nil, nil
	}

	database := clientMongo.Database(conf.Database)

	return nil, database, clientMongo
}

func ConnectDatabaseQueryError(conf *MongoConfig) (error, *mongo.Database, *mongo.Client) {

	readprefConfig := readpref.Primary()

	options := options.Client()
	options.ApplyURI(conf.Uri)
	options.SetReadPreference(readprefConfig)

	clientMongo, err := mongo.NewClient(options)
	if err != nil {
		return err, nil, nil
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = clientMongo.Connect(ctx)
	if err != nil {
		return err, nil, nil
	}

	database := clientMongo.Database(conf.Database)
	return nil, database, clientMongo
}
