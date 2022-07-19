package score

import (
	context "context"
	"party/models"
	"party/service"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type ScoreServiceImpl struct {
	UnimplementedScoreServiceServer
	Storage models.Storage
}

func (scoreService ScoreServiceImpl) UpdateScore(ctx context.Context, score *Score) (*emptypb.Empty, error) {
	partyService := service.PartyService{}
	return &emptypb.Empty{}, partyService.UpdateScore(scoreService.Storage, score.Partyid, score.Email, score.Score)
}
