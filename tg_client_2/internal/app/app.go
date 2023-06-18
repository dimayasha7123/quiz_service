package app

import (
	"github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/commands"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/queries"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
)

type Queries struct {
	EnqQuizHandler   queries.EndQuizHandler
	CurrQuestHandler queries.CurrentQuestHandler
	QuizzesHandler   queries.QuizzesHandler
	TopByQuizHandler queries.TopByQuizHandler
	StartHandler     queries.StartHandler
	StartQuizHandler queries.StartQuizHandler
}

type Commands struct {
	BreakHandler   commands.BreakHandler
	ConfirmHandler commands.ConfirmHandler
	SwitchHandler  commands.SwitchHandler
}

type App struct {
	Queries  Queries
	Commands Commands
}

func New(sessions domain.Sessions, quizClient api.QuizServiceClient) App {
	return App{
		Queries: Queries{
			EnqQuizHandler:   queries.NewEndQuizHandler(sessions, quizClient),
			CurrQuestHandler: queries.NewCurrentQuestHandler(sessions),
			QuizzesHandler:   queries.NewQuizzesHandler(quizClient),
			TopByQuizHandler: queries.NewTopByQuizHandler(sessions, quizClient),
			StartHandler:     queries.NewStartHandler(sessions, quizClient),
			StartQuizHandler: queries.NewStartQuizHandler(sessions, quizClient),
		},
		Commands: Commands{
			BreakHandler:   commands.NewBreakHandler(sessions),
			ConfirmHandler: commands.NewConfirmHandler(sessions, quizClient),
			SwitchHandler:  commands.NewSwitchHandler(sessions),
		},
	}
}
