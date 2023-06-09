package queries

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
)

type EnqQuizReq struct {
	UserID int64
}

type EnqQuizResp struct {
	Results app.Results
}

type endQuizHandler struct {
	sessions   domain.Sessions
	quizClient api.QuizServiceClient
}

func NewEndQuizHandler(sessions domain.Sessions, quizClient api.QuizServiceClient) endQuizHandler {
	return endQuizHandler{
		sessions:   sessions,
		quizClient: quizClient,
	}
}

func (h endQuizHandler) Handle(ctx context.Context, req EnqQuizReq) (EnqQuizResp, error) {
	err := h.sessions.EndQuizForUser(ctx, req.UserID)
	if err != nil {
		return EnqQuizResp{}, fmt.Errorf("can't end quiz for user with id = %v: %v", req.UserID, err)
	}

	user, err := h.sessions.GetUserByID(ctx, req.UserID)

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

	results := convertResultsFromApiToApp(ctx, h.sessions, qcResp)

	return EnqQuizResp{Results: results}, nil
}
