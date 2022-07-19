package service

import (
	"context"
	"encoding/json"
	"log"
	"party/models"
	partynamerepomongo "party/repository/partyNameRepoMongo"
	"party/repository/partyRepoSql"
	partyscorerepomongo "party/repository/partyScoreRepoMongo"
	"party/repository/userRepoMongo"
	"sync"

	"github.com/go-redis/redis/v9"
)

type PartyService struct{}

var ctx = context.Background()

func (ps *PartyService) GetLeaderboard(storage models.Storage, partyid string) (models.Leaderboard, error) {
	cache := storage.RedisCache
	result, err := cache.Get(ctx, partyid).Result()
	if err == nil {
		leaderboard := models.Leaderboard{}
		json.Unmarshal([]byte(result), &leaderboard)
		return leaderboard, nil
	}

	postgresDB := storage.PostgresPartyDB
	partyRepo := partyRepoSql.PartyRepo{}
	emails, err := partyRepo.GetPartyMemberIDs(postgresDB, partyid)
	if err != nil {
		log.Default().Panic(err)
		return models.Leaderboard{}, err
	}
	// find all usernames based on emails from user db (mongoDB)
	var wg sync.WaitGroup
	wg.Add(3)
	var emailUsernames map[string]string
	go func() {
		userRepo := userRepoMongo.UserRepoMongo{}
		emailUsernames = userRepo.GetUsernamesByEmails(storage.MongoUserDB, emails)
		wg.Done()
	}()

	// find party name
	var partyName string
	go func() {
		partyNameRepo := partynamerepomongo.PartyNameRepoMongo{}
		partyName, _ = partyNameRepo.GetPartyName(storage.MongoPartyDB, partyid)
		wg.Done()
	}()

	// find score for each email
	var emailScore map[string]int
	go func() {
		scoreRepo := partyscorerepomongo.PartyScoreRepoMongo{}
		emailScore = scoreRepo.GetScoresByIDs(storage.MongoPartyDB, partyid, emails)
		wg.Done()
	}()

	wg.Wait()

	// construct leaderboard
	members := make([]models.Member, len(emailScore))
	for _, email := range emails {
		member := models.Member{
			Email:    email,
			Username: emailUsernames[email],
			Score:    emailScore[email],
		}
		members = append(members, member)
	}
	leaderboard := models.Leaderboard{
		Name:    partyName,
		Members: members,
	}

	// set leaderboard into redis
	leaderboardByte, err := json.Marshal(leaderboard)
	if err != nil {
		log.Default().Panic(err)
		return leaderboard, err
	}
	cache.Set(ctx, partyid, leaderboardByte, 0)

	// return leaderboard
	return leaderboard, nil
}

func (ps *PartyService) UpdateScore(storage models.Storage, partyid string, email string, score int32) error {
	// policy: write around cache
	// remove leaderboard entry from cache
	ps.removePartyFromCache(storage.RedisCache, partyid)

	// update score db with new score
	scoreRepo := partyscorerepomongo.PartyScoreRepoMongo{}
	err := scoreRepo.UpdateScore(storage.MongoPartyDB, partyid, email, score)
	if err != nil {
		log.Default().Panic(err)
	}
	return err
}

func (ps *PartyService) AddMemberToParty(storage models.Storage, partyid string, email string) error {
	// remove from cache
	ps.removePartyFromCache(storage.RedisCache, partyid)

	// add record to party (member) db
	memberRepo := partyRepoSql.PartyRepo{}
	err := memberRepo.AddMemberToParty(storage.PostgresPartyDB, partyid, email)
	if err != nil {
		return err
	}

	// add record to score db
	scoreRepo := partyscorerepomongo.PartyScoreRepoMongo{}
	return scoreRepo.AddScore(storage.MongoPartyDB, partyid, email, 0)
}

func (ps *PartyService) AddParty(storage models.Storage, name string) (interface{}, error) {
	// add new party to party name db
	partyNameRepo := partynamerepomongo.PartyNameRepoMongo{}
	return partyNameRepo.AddParty(storage.MongoPartyDB, name)
}

func (ps *PartyService) removePartyFromCache(cache *redis.Client, partyid string) {
	err := cache.Del(ctx, partyid).Err()
	if err != nil {
		log.Default().Panic(err)
	}
}
