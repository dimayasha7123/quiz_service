package app

import (
	"context"
	pb "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q *qserver) SendAnswers(ctx context.Context, req *pb.AnswersPack) (*pb.SingleTop, error) {
	completed, err := q.repo.CheckPartyCompleted(ctx, req.QuizPartyID)
	if err != nil {
		return nil, status.Error(codes.Internal, "error when check party")
	}
	if completed {
		return nil, status.Error(codes.AlreadyExists, "party already completed")
	}

	rightAnswers, err := q.repo.GetRightAnswers(ctx, req.QuizPartyID)
	if err != nil {
		return nil, status.Error(codes.Internal, "error when getting right answers")
	}

	if len(req.Answers) != len(rightAnswers) {
		return nil, status.Error(codes.InvalidArgument, "wrong number of answers")
	}

	var points int32

	for i, a := range req.Answers {
		if len(a.RightAnswerNumbers) != len(rightAnswers[i]) {
			continue
		}
		hasWrong := false
		for j, an := range a.RightAnswerNumbers {
			if an != rightAnswers[i][j] {
				hasWrong = true
				break
			}
		}
		if hasWrong {
			continue
		}
		points++
	}

	singleTop, err := q.repo.CompleteParty(ctx, req.QuizPartyID, points)
	if err != nil {
		return nil, status.Error(codes.Internal, "error when writing results")
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
