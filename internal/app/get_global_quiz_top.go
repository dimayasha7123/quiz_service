package app

import (
	"context"
	pb "github.com/dimayasha7123/quiz_service/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q *qserver) GetGlobalQuizTop(ctx context.Context, req *pb.QuizID) (*pb.GlobalTop, error) {
	globalTop, err := q.repo.GetQuizTop(ctx, req.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, "error when counting of votes")
	}

	quizTopResult := make([]*pb.PartyResults, len(globalTop.Results))

	for i, r := range globalTop.Results {
		quizTopResult[i] = &pb.PartyResults{
			Name:       r.Name,
			PointCount: r.Points,
			Place:      r.Place,
		}
	}

	return &pb.GlobalTop{Results: quizTopResult}, nil
}
