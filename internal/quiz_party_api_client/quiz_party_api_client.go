package quizApi

import (
	"net/http"
	"time"
)

const (
	maxAnsCount = 6
	questCount  = 10
)

type QuizPartyApiClient struct {
	apiKey string
	cl     http.Client
}

func New(apiKey string) *QuizPartyApiClient {
	return &QuizPartyApiClient{
		cl:     http.Client{Timeout: 10 * time.Second},
		apiKey: apiKey,
	}
}
