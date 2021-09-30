package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func isInList(val string, slice []string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func WriteDB(stationc chan map[string]string, donec chan string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Println(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	collection := client.Database("OpenData").Collection("stations")
	don := []string{}
	for m := range stationc {
		_, isIn := isInList(m["_id"], don)
		if !isIn {
			b, err := bson.Marshal(m)
			if err != nil {
				log.Println("Cannot marshal map.")
			}
			o := options.Replace()
			o.SetUpsert(true)
			_, err = collection.ReplaceOne(
				context.Background(),
				bson.D{primitive.E{Key: "_id", Value: m["_id"]}},
				b,
				o,
			)
			if err != nil {
				log.Printf("%+v\n", err)
			}
		}
		don = append(don, m["_id"])

	}
	donec <- "Mongo is done.\n"
}
