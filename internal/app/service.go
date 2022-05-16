package app

import (
	pb "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/pkg/api"
)

type qserver struct {
	repo Repository
	pb.UnimplementedQuizServiceServer
}

func New(repo Repository) *qserver {
	return &qserver{repo: repo}
}
