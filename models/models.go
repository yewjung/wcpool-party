package models

import (
	"database/sql"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Leaderboard struct {
	Name    string   `json:"name"`
	Members []Member `json:"members"`
}

type Member struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Score    int    `json:"score"`
}

type Error struct {
	Message string `json:"message"`
}

type Storage struct {
	PostgresPartyDB *sql.DB
	MongoUserDB     *mongo.Client
	MongoPartyDB    *mongo.Client
	RedisCache      *redis.Client
}

type EmailUsername struct {
	Email    string `bson:"_id"`
	Username string `bson:"username"`
}

type EmailScore struct {
	Email string `bson:"_id"`
	Score int    `bson:"score"`
}
type MemberScore struct {
	PartyId string `json:"partyid"`
	Email   string `json:"email"`
	Score   int32  `json:"score"`
}

type Party struct {
	PartyId string `json:"partyid"`
	Name    string `json:"name"`
}
