package app

import (
	pb "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/pkg/api"
)

type qserver struct {
	repo     Repository
	qPartyCl QuizPartyApiClient
	pb.UnimplementedQuizServiceServer
}

func New(repo Repository, qPartyClient QuizPartyApiClient) *qserver {
	return &qserver{repo: repo, qPartyCl: qPartyClient}
}
