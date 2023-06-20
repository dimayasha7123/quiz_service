package queries

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/models"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
)

type EnqQuizReq struct {
	UserID int64
}

type EnqQuizResp struct {
	Results models.Results
}

type EndQuizHandler struct {
	sessions   domain.Sessions
	quizClient api.QuizServiceClient
}

func NewEndQuizHandler(sessions domain.Sessions, quizClient api.QuizServiceClient) EndQuizHandler {
	return EndQuizHandler{
		sessions:   sessions,
		quizClient: quizClient,
	}
}

func (h EndQuizHandler) Handle(ctx context.Context, req EnqQuizReq) (EnqQuizResp, error) {
	err := h.sessions.EndQuizForUser(ctx, req.UserID)
	if err != nil {
		return EnqQuizResp{}, fmt.Errorf("can't end quiz for user with id = %v: %v", req.UserID, err)
	}

	user, err := h.sessions.UserByID(ctx, req.UserID)
	if err != nil {
		return EnqQuizResp{}, fmt.Errorf("can't get user by id = %v: %v", req.UserID, err)
	}

	answers := make([]*api.QuestionRightAnswers, 0, 10)
	for _, quest := range user.Party.Questions {
		ansNums := make([]int32, 0, 6)
		for num, answer := range quest.Answers {
			if answer.Picked {
				ansNums = append(ansNums, int32(num))
			}
		}
		answers = append(answers, &api.QuestionRightAnswers{
			RightAnswerNumbers: ansNums,
		})
	}

	qcResp, err := h.quizClient.SendAnswers(ctx,
		&api.AnswersPack{
			QuizPartyID: user.Party.ID,
			Answers:     answers,
		},
	)
	if err != nil {
		return EnqQuizResp{}, fmt.Errorf("can't send answers to quiz server with quizPartyID = %v: %v", user.Party.ID, err)
	}

	results := convertResultsFromApiToApp(ctx, h.sessions, qcResp)

	return EnqQuizResp{Results: results}, nil
}
