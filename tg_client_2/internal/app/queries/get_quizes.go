package queries

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app"
	"google.golang.org/protobuf/types/known/emptypb"
)

type QuizzesResp struct {
	QuizList app.QuizList
}

type QuizzesHandler struct {
	quizClient api.QuizServiceClient
}

func NewQuizzesHandler(quizClient api.QuizServiceClient) QuizzesHandler {
	return QuizzesHandler{
		quizClient: quizClient,
	}
}

func (h QuizzesHandler) Handle(ctx context.Context) (QuizzesResp, error) {
	qcResp, err := h.quizClient.GetQuizList(ctx, &emptypb.Empty{})
	if err != nil {
		return QuizzesResp{}, fmt.Errorf("can't get quizzes from quiz service")
	}

	quizList := make(app.QuizList, 0, len(qcResp.QList))
	for _, quiz := range qcResp.QList {
		quizList = append(quizList, app.Quiz{
			ID:    quiz.ID,
			Title: quiz.Name,
		})
	}

	return QuizzesResp{
		QuizList: quizList,
	}, nil
}
