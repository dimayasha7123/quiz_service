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
	"testing"
)

func TestGetQuizTop(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	mockRepo := app.NewRepositoryMock(mc)

	sTop := m.SingleTop{
		UserResults: m.PartyResults{
			Name:   "Petya",
			Points: 8,
			Place:  2,
		},
		GlobalTop: m.GlobalTop{Results: []m.PartyResults{
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
		}},
	}

	mockRepo.GetUserQuizTopMock.Return(sTop, nil)

	svc := app.New(mockRepo, NewTestApiClient())
	ctx := context.Background()

	resp, err := svc.GetQuizTop(ctx, &pb.QuizUserInfo{
		UserID: 3,
		QuizID: 4,
	})
	assert.Nil(t, err)

	expectedSTop := pb.SingleTop{
		UserResults: &pb.PartyResults{
			Name:       "Petya",
			PointCount: 8,
			Place:      2,
		},
		QuizTop: &pb.GlobalTop{Results: []*pb.PartyResults{
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
		}},
	}
	assert.Equal(t, expectedSTop, *resp)

	mockRepo.GetUserQuizTopMock.Return(m.SingleTop{}, fmt.Errorf("some db error"))
	resp, err = svc.GetQuizTop(ctx, &pb.QuizUserInfo{
		UserID: 6,
		QuizID: 978,
	})
	assert.Equal(t, status.Error(codes.Internal, "error when counting of votes"), err)
	assert.Nil(t, resp)

}
