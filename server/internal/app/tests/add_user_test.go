package tests

import (
	"context"
	"fmt"
	app2 "github.com/dimayasha7123/quiz_service/server/internal/app"
	pb "github.com/dimayasha7123/quiz_service/server/pkg/api"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAddUser(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	ctx := context.Background()
	mockRepo := app2.NewRepositoryMock(mc)
	svc := app2.New(mockRepo, NewTestApiClient())

	t.Run("success", func(t *testing.T) {
		mockRepo.AddUserMock.Return(3, nil)
		name := "Dimyasha"

		resp, err := svc.AddUser(ctx, &pb.User{Name: name})

		assert.NoError(t, err)
		assert.EqualValues(t, 3, resp.ID)
	})
	t.Run("user already exists", func(t *testing.T) {
		mockRepo.AddUserMock.Return(-1, fmt.Errorf("not found"))
		name := "Dimyasha"

		resp, err := svc.AddUser(ctx, &pb.User{Name: name})

		assert.EqualError(t, status.Error(codes.AlreadyExists, fmt.Sprintf("user with name <%s> is already exists", name)), err.Error())
		assert.EqualValues(t, -1, resp.ID)
	})
}
