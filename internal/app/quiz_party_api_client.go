package app

import "github.com/dimayasha7123/quiz_service/internal/models"

type QuizPartyApiClient interface {
	GetParty(tag string) (*models.Party, error)
}
