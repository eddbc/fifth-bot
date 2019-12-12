package db

import (
	"context"
	"encoding/binary"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var db *mongo.Client

func DBInit() {
	// Set client options
	// TODO get DB credentials from ENV
	clientOptions := options.Client().ApplyURI("mongodb://admin:password@localhost:27017")
	// Connect to MongoDB
	cl, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = cl.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	db = cl
	log.Print("Connected to MongoDB")

	_, err = Db().Collection("tokens").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{{ "characterid",1, }},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Print(err)
	}

	_, err = Db().Collection("discord").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{{"discordid", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Print(err)
	}
	_, err = Db().Collection("discord").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{{ "characterid",1, }},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Print(err)
	}

}

func Db() *mongo.Database {
	return db.Database("FleetBot")
}

// Itob returns an 8-byte big endian representation of v.
func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
