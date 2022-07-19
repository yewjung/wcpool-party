package controller

import (
	"context"
	"net/http"
	"party/authorization"
	"party/models"
	"party/service"
	"party/utils"

	"github.com/gorilla/mux"
)

type PartyController struct {
	AuthClient authorization.AuthorizationClient
}

func (pc *PartyController) GetLeaderboard(storage models.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		partyid := mux.Vars(r)["partyid"]

		// TODO: need to verify that user has access to this party

		// Looking for a valid feature
		if !pc.checkAuthorization(w, r, partyid, []authorization.Option{authorization.Option_PARTY_ID}) {
			return
		}
		partyService := service.PartyService{}
		leaderboard, err := partyService.GetLeaderboard(storage, partyid)
		utils.HandleResponse(w, err, leaderboard)
	}
}

// this method should be called through grpc
func (pc *PartyController) UpdateScore(storage models.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto := utils.GetReqBody(r, models.MemberScore{})
		partyService := service.PartyService{}
		err := partyService.UpdateScore(storage, dto.PartyId, dto.Email, dto.Score)
		utils.HandleResponse(w, err, nil)
	}
}

func (pc *PartyController) AddMemberToParty(storage models.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto := utils.GetReqBody(r, models.MemberScore{})
		// TODO: verify permission -> only admin can approve member
		// permission should be a string like :
		// <partyid>$admin
		// use email to find list of permission, use partyid to construct <partyid>$admin
		// user's profile can have a field called permissions = [..., ..., ...]
		if !pc.checkAuthorization(w, r, dto.PartyId, []authorization.Option{authorization.Option_PARTY_ID}) {
			return
		}
		partyService := service.PartyService{}
		err := partyService.AddMemberToParty(storage, dto.PartyId, dto.Email)
		utils.HandleResponse(w, err, nil)
	}
}
func (pc *PartyController) AddParty(storage models.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		party := utils.GetReqBody(r, models.Party{})
		// TODO: validate token
		pc.checkAuthorization(w, r, party.PartyId, nil)
		partyService := service.PartyService{}
		result, err := partyService.AddParty(storage, party.Name)
		utils.HandleResponse(w, err, result)
	}
}

func (pc *PartyController) checkAuthorization(w http.ResponseWriter, r *http.Request, partyid string, options []authorization.Option) bool {
	boolWrap, err := pc.AuthClient.VerifyPartyID(context.Background(), &authorization.Verification{
		Token:   r.Header.Get("Authorization"),
		Partyid: partyid,
		Options: options,
	})
	if err != nil {
		utils.HandleResponse(w, err, nil)
		return false
	}
	if !boolWrap.Value {
		utils.SendError(w, http.StatusUnauthorized, models.Error{
			Message: "Unauthorized action",
		})
		return false
	}
	return true
}
