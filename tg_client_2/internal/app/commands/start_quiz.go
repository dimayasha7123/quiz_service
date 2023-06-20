package commands

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/models"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain/states"
)

type StartQuizReq struct {
	UserQuizIDs models.UserQuizIDs
}

type StartQuizHandler struct {
	sessions   domain.Sessions
	quizClient api.QuizServiceClient
}

func NewStartQuizHandler(sessions domain.Sessions, quizClient api.QuizServiceClient) StartQuizHandler {
	return StartQuizHandler{
		sessions:   sessions,
		quizClient: quizClient,
	}
}

func (h StartQuizHandler) Handle(ctx context.Context, req StartQuizReq) error {
	user, err := h.sessions.UserByID(ctx, req.UserQuizIDs.UserID)
	if err != nil {
		return fmt.Errorf("can't get user from sessions: %v", err)
	}

	state, err := h.sessions.UserState(ctx, req.UserQuizIDs.UserID)
	if err != nil {
		return fmt.Errorf("can't get user from sessions: %v", err)
	}
	if state == states.Quiz {
		return fmt.Errorf("user now in party, can't start new before finishing previous")
	}

	qcResp, err := h.quizClient.StartQuizParty(ctx, &api.QuizUserInfo{
		UserID: user.QuizServiceID,
		QuizID: req.UserQuizIDs.QuizID,
	})
	if err != nil {
		return fmt.Errorf("can't start new quiz party via quiz service")
	}

	questions := make(domain.NewQuestions, 0, len(qcResp.Questions))
	for _, question := range qcResp.Questions {
		answers := make(domain.NewAnswers, 0, len(question.AnswerOptions))
		for _, answer := range question.AnswerOptions {
			answers = append(answers, domain.NewAnswer(answer))
		}
		questions = append(questions, domain.NewQuestion{
			Title:   question.Title,
			Answers: answers,
		})
	}
	party := domain.NewParty{
		ID:        qcResp.QuizPartyID,
		Questions: questions,
	}

	err = h.sessions.StartNewQuizForUser(ctx, req.UserQuizIDs.UserID, party)
	if err != nil {
		return fmt.Errorf("can't start new quiz party for user: %v", err)
	}

	return nil
}
