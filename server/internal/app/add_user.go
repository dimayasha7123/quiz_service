package app

import (
	"context"
	"fmt"
	pb "github.com/dimayasha7123/quiz_service/server/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q *qserver) AddUser(ctx context.Context, req *pb.User) (*pb.UserID, error) {
	userID, err := q.repo.GetUser(ctx, req.Name)
	if err != nil {
		return &pb.UserID{ID: -1}, status.Error(
			codes.Internal,
			fmt.Sprintf("can't check if user <%s> already exists", req.Name),
		)
	}
	if userID != -1 {
		return &pb.UserID{ID: userID}, nil
	}

	userID, err = q.repo.AddUser(ctx, req.Name)
	if err != nil {
		return &pb.UserID{ID: -1}, status.Error(
			codes.Internal,
			fmt.Sprintf("can't add user <%s> to database", req.Name),
		)
	}

	return &pb.UserID{ID: userID}, nil
}
