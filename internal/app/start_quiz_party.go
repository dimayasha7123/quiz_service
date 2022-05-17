package app

import (
	"context"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/models"
	pb "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (q *qserver) StartQuizParty(ctx context.Context, req *pb.QuizUserInfo) (*pb.QuizParty, error) {

	// + создаем Party с указанными QuizID и UserAccountID
	// + получаем вопросы по API с использованием QuizTitle по QuizID
	// + пытаемся добавить пак вопросов в БД, чтобы иметь их у себя
	// идем по всем вопросам:
	// +   добавляем вопрос в слайс, чтобы отправить пользователю
	// -   добавляем новую запись в partyQuestion

	// возможно тут все стоит все завернуть в транзакцию, ибо тут 4 обращения к репозу, но пока пусть будет так

	party := models.Party{
		ID:            0,
		UserAccountID: req.UserID,
		QuizID:        req.QuizID,
		Completed:     false,
		Points:        0,
		Questions:     nil,
	}

	partyID, err := q.repo.AddParty(ctx, party)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "no such UserAccountID or QuizID")
	}
	party.ID = partyID

	quiz, err := q.repo.GetQuiz(ctx, party.QuizID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "no such QuizID")
	}

	apiParty, err := q.qPartyCl.GetParty(quiz.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, "quizApiClient error")
	}

	party.Questions = apiParty.Questions

	countAddedQuestions, err := q.repo.AddQuestionsIfNot(ctx, party.Questions, party.QuizID)
	if err != nil {
		return nil, status.Error(codes.Internal, "error when add questions")
	}
	log.Printf("Add %d quest.", countAddedQuestions)

	err = q.repo.AddAllPartyQuestion(ctx, party.Questions, party.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, "error when add questions to party")
	}

	pbQuestions := make([]*pb.Question, len(party.Questions))

	for i, qst := range party.Questions {

		answerOptions := make([]string, len(qst.Answers))
		for j, ans := range qst.Answers {
			answerOptions[j] = ans.Title
		}

		pbQuestions[i] = &pb.Question{
			Title:         qst.Title,
			AnswerOptions: answerOptions,
		}
	}

	return &pb.QuizParty{QuizPartyID: party.ID, Questions: pbQuestions}, nil
}
