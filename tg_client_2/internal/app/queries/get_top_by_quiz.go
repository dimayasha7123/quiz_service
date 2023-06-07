package queries

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
	"github.com/dimayasha7123/quiz_service/utils/logger"
	"strconv"
)

type GetTopByQuizReq struct {
	UserQuizIDs app.UserQuizIDs
}

type GetTopByQuizResp struct {
	UserResults *app.UserResults
	TopResults  app.TopResults
}

type getTopByQuizHandler struct {
	sessions   domain.Sessions
	quizClient api.QuizServiceClient
}

func NewGetTopByQuizHandler(sessions domain.Sessions, quizClient api.QuizServiceClient) getTopByQuizHandler {
	return getTopByQuizHandler{
		sessions:   sessions,
		quizClient: quizClient,
	}
}

func (h getTopByQuizHandler) Handle(ctx context.Context, req GetTopByQuizReq) (GetTopByQuizResp, error) {
	user, err := h.sessions.GetUserByID(ctx, req.UserQuizIDs.UserID)
	if err != nil {
		return GetTopByQuizResp{}, fmt.Errorf("can't get user from sessions: %v", err)
	}

	qcResp, err := h.quizClient.GetQuizTop(ctx, &api.QuizUserInfo{
		UserID: user.QuizServiceID,
		QuizID: req.UserQuizIDs.QuizID,
	})

	ret := GetTopByQuizResp{}
	if qcResp.UserResults != nil && qcResp.UserResults.Place != 0 && qcResp.UserResults.PointCount != 0 {
		ret.UserResults = &app.UserResults{
			Place:  qcResp.UserResults.Place,
			Points: int64(qcResp.UserResults.PointCount),
		}
	}

	if len(qcResp.QuizTop.Results) != 0 {
		topResults := make(app.TopResults, 0, len(qcResp.QuizTop.Results))

		for _, result := range qcResp.QuizTop.Results {
			topResults = append(topResults, app.ResultRow{
				Username: h.getName(ctx, result.Name),
				Points:   result.Place,
			})
		}

		ret.TopResults = topResults
	}

	return ret, nil
}

func (h getTopByQuizHandler) getName(ctx context.Context, name string) string {
	id, err := strconv.ParseInt(name, 10, 64)
	if err != nil {
		logger.Log.Errorf("can't convert ID = %v from string to int: %v", name, err)
		return name
	}

	user, err := h.sessions.GetUserByID(ctx, id)
	if err != nil {
		logger.Log.Errorf("can't get user by id = %v from sessions: %v", id, err)
		return name
	}

	return user.Name
}
