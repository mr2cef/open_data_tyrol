package mongo

import (
	"context"
	"log"
	"os"
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

func mongoGetColletion() *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_DB_HOST")))
	if err != nil {
		log.Println(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	c := client.Database(os.Getenv("MONGO_DB_DB")).Collection(os.Getenv("MONGO_DB_COLLECTION"))
	return c
}

func mongoWriteToCollection(c *mongo.Collection, m map[string]string) {
	b, err := bson.Marshal(m)
	if err != nil {
		log.Println("Cannot marshal map.")
	}
	o := options.Replace()
	o.SetUpsert(true)
	_, err = c.ReplaceOne(
		context.Background(),
		bson.D{primitive.E{Key: "_id", Value: m["_id"]}},
		b,
		o,
	)
	if err != nil {
		log.Printf("Cannot write to DB: %+v\n", err)
	}
}

func WriteDB(stationc chan map[string]string, donec chan string) {

	c := mongoGetColletion()
	defer func() {
		if err := c.Database().Client().Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()
	don := []string{}
	for m := range stationc {
		_, isIn := isInList(m["_id"], don)
		if !isIn {
			mongoWriteToCollection(c, m)
			don = append(don, m["_id"])
		}
	}
	donec <- "Mongo is done.\n"
}
