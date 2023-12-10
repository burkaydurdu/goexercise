package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//go:generate go vet ./...

func main() {
	var ctx = context.TODO()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo-master:27017,mongo-slave-1:27018/?replicaSet=rs0"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Connection Successfully", err)
	}

	databaseName := "goexercise"
	collectionName := "user"

	session, err := client.StartSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.EndSession(ctx)

	err = session.StartTransaction()
	if err != nil {
		log.Fatal(err)
	}

	ctx = mongo.NewSessionContext(ctx, session)

	collection := client.Database(databaseName).Collection(collectionName)
	_, err = collection.InsertOne(ctx, bson.M{"key": "value1"})
	if err != nil {
		_ = session.AbortTransaction(ctx)
		log.Fatal(err)
	}

	_, err = collection.UpdateOne(ctx, bson.M{"key": "value1"}, bson.M{"$set": bson.M{"key": "value2"}})
	if err != nil {
		_ = session.AbortTransaction(ctx)
		log.Fatal(err)
	}

	_ = session.AbortTransaction(ctx)

	// _ = session.CommitTransaction(ctx)

	fmt.Println("Transaction has canceled successfully")
}
