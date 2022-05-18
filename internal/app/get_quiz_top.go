package app

import (
	"context"
	pb "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q *qserver) GetQuizTop(ctx context.Context, req *pb.QuizUserInfo) (*pb.SingleTop, error) {
	singleTop, err := q.repo.GetUserQuizTop(ctx, req.QuizID, req.UserID)
	if err != nil {
		return nil, status.Error(codes.Internal, "error when counting of votes")
	}

	quizTopResult := make([]*pb.PartyResults, len(singleTop.GlobalTop.Results))

	for i, r := range singleTop.GlobalTop.Results {
		quizTopResult[i] = &pb.PartyResults{
			Name:       r.Name,
			PointCount: r.Points,
			Place:      r.Place,
		}
	}

	pbSingleTop := pb.SingleTop{
		UserResults: &pb.PartyResults{
			Name:       singleTop.UserResults.Name,
			PointCount: singleTop.UserResults.Points,
			Place:      singleTop.UserResults.Place,
		},
		QuizTop: &pb.GlobalTop{Results: quizTopResult},
	}

	return &pbSingleTop, nil
}
