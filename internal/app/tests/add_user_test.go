package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/dimayasha7123/quiz_service/internal/app"
	pb "github.com/dimayasha7123/quiz_service/pkg/api"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAddUser(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	mockRepo := app.NewRepositoryMock(mc)
	mockRepo.AddUserMock.Return(3, nil)

	svc := app.New(mockRepo, NewTestApiClient())
	name := "Dimyasha"
	ctx := context.Background()

	resp, err := svc.AddUser(ctx, &pb.User{Name: name})
	assert.Nil(t, err)
	assert.Equal(t, int64(3), resp.ID)

	mockRepo.AddUserMock.Return(-1, fmt.Errorf("not found"))
	resp, err = svc.AddUser(ctx, &pb.User{Name: name})
	assert.Equal(t, status.Error(codes.AlreadyExists, fmt.Sprintf("user with name <%s> is already exists", name)), err)
	assert.Equal(t, int64(-1), resp.ID)
}
