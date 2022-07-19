package driver

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/go-redis/redis/v9"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectPostgresPartyDB() *sql.DB {
	pgurl := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"partyDB", 5432, "user", "mysecretpassword", "user")
	db, err := sql.Open("postgres", pgurl)
	logFatal(err)

	err = db.Ping()
	logFatal(err)

	return db
}

func ConnectMongoPartyDB() *mongo.Client {
	uri := "mongodb://root:example@partydb:27017/"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("MongoDB partydb successfully connected and pinged.")
	return client
}

func ConnectMongoUserDB() *mongo.Client {
	uri := "mongodb://root:example@userdb:27017/"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("MongoDB userdb successfully connected and pinged.")
	return client
}

func ConnectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}
