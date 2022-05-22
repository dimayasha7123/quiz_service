package tests

import (
	"context"
	"fmt"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/app"
	pb "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
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