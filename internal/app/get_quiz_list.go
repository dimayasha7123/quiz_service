package app

import (
	"context"
	pb "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (q *qserver) GetQuizList(ctx context.Context, req *emptypb.Empty) (*pb.QuizList, error) {

	quizes, err := q.repo.GetQuizList(ctx)

	if err != nil {
		return nil, status.Error(codes.Internal, "Some DB error")
	}

	qList := make([]*pb.Quiz, len(quizes))

	for i, qz := range quizes {
		qList[i] = &pb.Quiz{
			ID:   qz.ID,
			Name: qz.Title,
		}
	}

	return &pb.QuizList{QList: qList}, nil
}
