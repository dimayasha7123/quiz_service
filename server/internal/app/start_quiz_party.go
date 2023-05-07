package app

import (
	"context"
	"github.com/dimayasha7123/quiz_service/server/internal/models"
	pb "github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/utils/logger"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q *qserver) StartQuizParty(ctx context.Context, req *pb.QuizUserInfo) (*pb.QuizParty, error) {
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

	startTime := time.Now()
	apiParty, err := q.qPartyCl.GetParty(quiz.Title)
	if err != nil {
		return nil, status.Error(codes.Internal, "quizApiClient error")
	}

	logger.Log.Infof("Time for QuizAPI request: %v", time.Since(startTime))

	party.Questions = apiParty.Questions

	countAddedQuestions, err := q.repo.AddQuestionsIfNot(ctx, &party.Questions, party.QuizID)
	if err != nil {
		return nil, status.Error(codes.Internal, "error when add questions")
	}
	logger.Log.Infof("Add %d questions", countAddedQuestions)

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
