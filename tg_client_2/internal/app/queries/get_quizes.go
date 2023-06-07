package queries

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GetQuizzesResp struct {
	QuizList app.QuizList
}

type getQuizzesHandler struct {
	quizClient api.QuizServiceClient
}

func NewGetQuizzesHandler(quizClient api.QuizServiceClient) getQuizzesHandler {
	return getQuizzesHandler{
		quizClient: quizClient,
	}
}

func (h getQuizzesHandler) Handle(ctx context.Context) (GetQuizzesResp, error) {
	qcResp, err := h.quizClient.GetQuizList(ctx, &emptypb.Empty{})
	if err != nil {
		return GetQuizzesResp{}, fmt.Errorf("can't get quizzes from quiz service")
	}

	quizList := make(app.QuizList, 0, len(qcResp.QList))
	for _, quiz := range qcResp.QList {
		quizList = append(quizList, app.Quiz{
			ID:    quiz.ID,
			Title: quiz.Name,
		})
	}

	return GetQuizzesResp{
		QuizList: quizList,
	}, nil
}
