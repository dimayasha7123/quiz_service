package app

import (
	pb "github.com/dimayasha7123/quiz_service/server/pkg/api"
)

type qserver struct {
	repo     Repository
	qPartyCl QuizPartyApiClient
	pb.UnimplementedQuizServiceServer
}

func New(repo Repository, qPartyClient QuizPartyApiClient) *qserver {
	return &qserver{repo: repo, qPartyCl: qPartyClient}
}
