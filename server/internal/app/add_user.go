package app

import (
	"context"
	"fmt"
	pb "github.com/dimayasha7123/quiz_service/server/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q *qserver) AddUser(ctx context.Context, req *pb.User) (*pb.UserID, error) {

	userID, err := q.repo.AddUser(ctx, req.Name)

	if err != nil {
		return &pb.UserID{ID: -1}, status.Error(
			codes.AlreadyExists,
			fmt.Sprintf("user with name <%s> is already exists", req.Name),
		)
	}

	return &pb.UserID{ID: userID}, nil
}
