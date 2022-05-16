package app

import (
	"context"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/models"
	pb "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q *qserver) CreateQuizParty(ctx context.Context, req *pb.QuizUserInfo) (*pb.QuizPartyID, error) {

	party := models.Participation{
		ID:            0,
		UserAccountID: req.UserID,
		QuizID:        req.QuizID,
	}

	partyID, err := q.repo.AddParticipation(ctx, party)
	if err != nil {
		return &pb.QuizPartyID{ID: -1}, status.Error(codes.InvalidArgument, "no such UserAccountID or QuizID")
	}

	return &pb.QuizPartyID{ID: partyID}, nil
}
