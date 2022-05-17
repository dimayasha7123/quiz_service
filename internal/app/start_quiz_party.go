package app

import (
	"context"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/models"
	pb "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q *qserver) StartQuizParty(ctx context.Context, req *pb.QuizUserInfo) (*pb.QuizParty, error) {

	// COMPLETED создаем Quiz с указанными QuizID и UserAccountID
	// получаем вопросы по API с использованием QuizTitle по QuizID
	// идем по всем вопросам:
	//     добавляем вопрос в слайс, чтобы отправить пользователю
	//     если вопроса нет в БД (ищем по Question.Title), то добавляем его
	//     добавляем новую запись в partyQuestion

	party := models.Party{
		ID:            0,
		UserAccountID: req.UserID,
		QuizID:        req.QuizID,
		Completed:     false,
		Points:        0,
	}

	partyID, err := q.repo.AddParty(ctx, party)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "no such UserAccountID or QuizID")
	}

	party.ID = partyID

	return &pb.QuizParty{QuizPartyID: partyID}, nil
}
