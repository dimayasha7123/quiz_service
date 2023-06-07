package domain

import (
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain/states"
)

type User struct {
	TelegramID    int64
	QuizServiceID int64
	Name          string
	State         states.State
	Party         Party
}

type Party struct {
	ID              int64
	CurrentQuestion int64
	Questions       Questions
}

type Questions []Question

type Question struct {
	Title   string
	Answers Answers
}

type Answers []Answer

type Answer struct {
	Title  string
	Picked bool
}

// for new

type NewParty struct {
	ID        int64
	Questions NewQuestions
}

type NewQuestions []NewQuestion

type NewQuestion struct {
	Title   string
	Answers NewAnswers
}

type NewAnswers []NewAnswer

type NewAnswer string

func (u *User) inQuiz() bool {
	return u.State == states.Quiz
}

func (u *User) StartNewQuiz(newParty NewParty) error {
	if u.inQuiz() {
		return fmt.Errorf("can't start new quiz: user already in quiz")
	}

	questions := make(Questions, 0, len(newParty.Questions))
	for _, question := range newParty.Questions {
		answers := make(Answers, 0, len(question.Answers))
		for _, answer := range question.Answers {
			answers = append(answers, Answer{
				Title:  string(answer),
				Picked: false,
			})
		}
		questions = append(questions, Question{
			Title:   question.Title,
			Answers: answers,
		})
	}
	party := Party{
		ID:              newParty.ID,
		CurrentQuestion: 0,
		Questions:       questions,
	}

	u.Party = party
	u.State = states.Quiz

	return nil
}

func (u *User) GetCurrentQuestion() (bool, Question, error) {
	if !u.inQuiz() {
		return false, Question{}, fmt.Errorf("user not in quiz")
	}

	if len(u.Party.Questions) == 0 {
		return false, Question{}, fmt.Errorf("no questions in current party")
	}

	if u.Party.CurrentQuestion >= int64(len(u.Party.Questions)) {
		return false, Question{}, nil
	}

	return true, u.Party.Questions[u.Party.CurrentQuestion], nil
}
