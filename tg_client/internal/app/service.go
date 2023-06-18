package app

import (
	"net/http"
	"time"

	pb "github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/tg_client/internal/models"
)

type bclient struct {
	repo       repository
	apiKey     string
	quizClient pb.QuizServiceClient
	httpClient http.Client
	users      *models.SyncMap
}

func New(repo repository, apiKey string, quizClient pb.QuizServiceClient) *bclient {
	bc := bclient{
		repo:       repo,
		apiKey:     apiKey,
		quizClient: quizClient,
		httpClient: http.Client{Timeout: 30 * time.Second},
		users:      models.NewSyncMap(),
	}
	return &bc
}
