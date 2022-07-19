package main

import (
	"fmt"
	"log"
	"net/http"
	"party/authorization"
	"party/controller"
	"party/driver"
	"party/models"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func main() {
	partyDB := driver.ConnectPostgresPartyDB()
	mongoUserDB := driver.ConnectMongoUserDB()
	mongoPartyDB := driver.ConnectMongoPartyDB()
	redisCache := driver.ConnectRedis()
	storage := models.Storage{
		PostgresPartyDB: partyDB,
		MongoUserDB:     mongoUserDB,
		MongoPartyDB:    mongoPartyDB,
		RedisCache:      redisCache,
	}

	router := mux.NewRouter()

	partyController := controller.PartyController{AuthClient: getSecurityGrpcClient()}
	// /leaderboard/{partyid}/
	router.HandleFunc("/leaderboard/{partyid}", partyController.GetLeaderboard(storage)).Methods("GET")
	// /score data: {partyid, email, score}
	router.HandleFunc("/score", partyController.UpdateScore(storage)).Methods("POST")
	// /member data: {partyid, email}
	router.HandleFunc("/member", partyController.AddMemberToParty(storage)).Methods("POST")
	// /party data: {partyname}
	router.HandleFunc("/party", partyController.AddParty(storage)).Methods("POST")

	fmt.Println("Server is running at port 8090")

	log.Fatal(
		http.ListenAndServe(":8090",
			handlers.CORS(
				handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
				handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
				handlers.AllowedOrigins([]string{"*"}),
			)(router),
		),
	)

}

func getSecurityGrpcClient() authorization.AuthorizationClient {
	conn, err := grpc.Dial("security:8085")
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return authorization.NewAuthorizationClient(conn)
}
