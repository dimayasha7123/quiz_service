package tests

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/internal/app"
	m "github.com/dimayasha7123/quiz_service/internal/models"
	pb "github.com/dimayasha7123/quiz_service/pkg/api"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func TestGetQuizList(t *testing.T) {

	mc := minimock.NewController(t)
	defer mc.Finish()

	mockRepo := app.NewRepositoryMock(mc)
	mockRepo.GetQuizListMock.Return([]m.Quiz{
		{
			ID:    1,
			Title: "Linux",
		},
		{
			ID:    2,
			Title: "Bash",
		},
		{
			ID:    3,
			Title: "Docker",
		},
	}, nil)

	svc := app.New(mockRepo, NewTestApiClient())
	ctx := context.Background()

	resp, err := svc.GetQuizList(ctx, &emptypb.Empty{})
	assert.Nil(t, err)
	assert.Equal(t, pb.QuizList{QList: []*pb.Quiz{
		{
			ID:   1,
			Name: "Linux",
		},
		{
			ID:   2,
			Name: "Bash",
		},
		{
			ID:   3,
			Name: "Docker",
		},
	}}, *resp)

	mockRepo.GetQuizListMock.Return(nil, fmt.Errorf("some db error"))
	resp, err = svc.GetQuizList(ctx, &emptypb.Empty{})
	assert.Equal(t, status.Error(codes.Internal, "Some DB error"), err)
	assert.Nil(t, resp)
}
