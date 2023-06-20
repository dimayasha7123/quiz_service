package queries

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/models"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
)

type StartReq struct {
	UserInfo models.UserInfo
}

type StartResp struct {
	NewUser bool
}

type StartHandler struct {
	sessions   domain.Sessions
	quizClient api.QuizServiceClient
}

func NewStartHandler(sessions domain.Sessions, quizClient api.QuizServiceClient) StartHandler {
	return StartHandler{
		sessions:   sessions,
		quizClient: quizClient,
	}
}

func (h StartHandler) Handle(ctx context.Context, req StartReq) (StartResp, error) {
	exists := h.sessions.UserExistsByID(ctx, req.UserInfo.UserID)
	if exists {
		return StartResp{NewUser: false}, nil
	}

	username := fmt.Sprintf("%d", req.UserInfo.UserID)
	qcResp, err := h.quizClient.AddUser(ctx, &api.User{Name: username})
	if err != nil {
		return StartResp{}, fmt.Errorf("can't add user to quiz service: %v", err)
	}

	err = h.sessions.AddUser(ctx, req.UserInfo.UserID, qcResp.ID, req.UserInfo.UserName)
	if err != nil {
		return StartResp{}, fmt.Errorf("can't add user to sessions: %v", err)
	}

	return StartResp{NewUser: true}, nil
}
