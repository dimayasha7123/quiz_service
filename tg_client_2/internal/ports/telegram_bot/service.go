package telegram_bot

import (
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app"
	"net/http"
	"time"
)

type service struct {
	apiKey     string
	delay      time.Duration
	httpClient http.Client
	app        app.App
}

func New(apiKey string, app app.App, httpTimeout time.Duration, tgApiDelay time.Duration) service {
	return service{
		apiKey:     apiKey,
		delay:      tgApiDelay,
		httpClient: http.Client{Timeout: httpTimeout},
		app:        app,
	}
}
