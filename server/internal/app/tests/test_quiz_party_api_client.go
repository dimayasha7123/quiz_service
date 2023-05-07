package tests

import (
	"github.com/dimayasha7123/quiz_service/server/internal/quiz_party_api_client"
)

func NewTestApiClient() *quizApi.QuizPartyApiClient {
	return quizApi.New("GLAwwDtuyBMEwzWnQsq0Es4oAG8kMXHOfFdJsb1E")
}
