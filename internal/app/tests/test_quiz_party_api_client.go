package tests

import (
	quizApi "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/quiz_party_api_client"
)

func NewTestApiClient() *quizApi.QuizPartyApiClient {
	return quizApi.New("GLAwwDtuyBMEwzWnQsq0Es4oAG8kMXHOfFdJsb1E")
}
