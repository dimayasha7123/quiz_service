package app

import "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/models"

type QuizPartyApiClient interface {
	GetParty(tag string) (*models.Party, error)
}
