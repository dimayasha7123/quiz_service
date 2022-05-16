package app

import (
	"context"
	"fmt"
	pb "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q *qserver) AddUser(ctx context.Context, req *pb.User) (*pb.UserID, error) {
	// если есть юзер с таким именем, то возвращаем -1
	// иначе возвращаем id только что добавленного юзера

	userID, err := q.repo.AddUser(ctx, req.Name)

	if err != nil {
		return &pb.UserID{ID: -1}, status.Error(
			codes.AlreadyExists,
			fmt.Sprintf("user with name <%s> is already exists", req.Name),
		)
	}

	return &pb.UserID{ID: userID}, nil
}
