package app

import (
	"context"
	"errors"
	pb "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/pkg/api"
)

func (q *qserver) AddUser(ctx context.Context, req *pb.User) (*pb.UserID, error) {
	// если есть юзер с таким именем, то возвращаем -1
	// иначе возвращаем id только что добавленного юзера

	userID, err := q.repo.AddUser(ctx, req.Name)
	if err != nil {
		return &pb.UserID{ID: -1}, errors.New("user is already exists")
	}
	return &pb.UserID{ID: userID}, nil
}
