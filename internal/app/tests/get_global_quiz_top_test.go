package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/dimayasha7123/quiz_service/internal/app"
	m "github.com/dimayasha7123/quiz_service/internal/models"
	pb "github.com/dimayasha7123/quiz_service/pkg/api"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetGlobalQuizTop(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	mockRepo := app.NewRepositoryMock(mc)
	svc := app.New(mockRepo, NewTestApiClient())
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		gTop := m.GlobalTop{Results: []m.PartyResults{
			{
				Name:   "Vasya",
				Points: 10,
				Place:  1,
			},
			{
				Name:   "Petya",
				Points: 8,
				Place:  2,
			},
		}}
		mockRepo.GetQuizTopMock.Return(gTop, nil)
		expectedGTop := pb.GlobalTop{Results: []*pb.PartyResults{
			{
				Name:       "Vasya",
				PointCount: 10,
				Place:      1,
			},
			{
				Name:       "Petya",
				PointCount: 8,
				Place:      2,
			},
		}}

		resp, err := svc.GetGlobalQuizTop(ctx, &pb.QuizID{ID: 4})

		assert.NoError(t, err)
		assert.Equal(t, expectedGTop, *resp)
	})
	t.Run("with db error", func(t *testing.T) {
		mockRepo.GetQuizTopMock.Return(m.GlobalTop{}, fmt.Errorf("some db error"))

		resp, err := svc.GetGlobalQuizTop(ctx, &pb.QuizID{ID: 4})

		assert.EqualError(t, status.Error(codes.Internal, "error when counting of votes"), err.Error())
		assert.Nil(t, resp)
	})

}
