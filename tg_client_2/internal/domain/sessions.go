package domain

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain/states"
	"sync"
)

type Sessions struct {
	mutex *sync.RWMutex
	repo  repository
	data  map[int64]User
}

func NewSessions(ctx context.Context, repo repository) (Sessions, error) {
	users, err := repo.GetAllUsers(ctx)
	if err != nil {
		return Sessions{}, fmt.Errorf("can't get users from repo: %v", err)
	}

	data := make(map[int64]User)
	doubled := make([]int64, 0)
	for _, user := range users {
		_, ok := data[user.ID]
		if ok {
			doubled = append(doubled, user.ID)
			continue
		}
		data[user.ID] = User{
			TelegramID:    user.ID,
			QuizServiceID: user.QSID,
			Name:          user.Name,
			State:         states.Lobby,
		}
	}

	if len(doubled) != 0 {
		return Sessions{}, fmt.Errorf("can't create data storage, find doubled ids: %v", doubled)
	}

	return Sessions{
		mutex: &sync.RWMutex{},
		repo:  repo,
		data:  data,
	}, nil
}

func (s *Sessions) CheckUserExistsByID(ctx context.Context, id int64) bool {
	s.mutex.RLock()
	_, ok := s.data[id]
	s.mutex.RUnlock()

	return ok
}

func (s *Sessions) GetUserByID(ctx context.Context, id int64) (User, error) {
	s.mutex.RLock()
	user, ok := s.data[id]
	s.mutex.RUnlock()

	if ok {
		return user, nil
	}

	return User{}, fmt.Errorf("no such user")
}

func (s *Sessions) AddUser(ctx context.Context, id int64, qsid int64, name string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.data[id]

	if ok {
		return fmt.Errorf("user with id = %v already exists", id)
	}

	err := s.repo.AddUser(ctx, RepoUser{
		ID:   id,
		QSID: qsid,
		Name: name,
	})

	if err != nil {
		return fmt.Errorf("can't add user to repo: %v", err)
	}

	user := User{
		TelegramID:    id,
		QuizServiceID: qsid,
		Name:          name,
		State:         states.Lobby,
	}

	s.data[user.TelegramID] = user

	return nil
}

func (s *Sessions) GetUserState(ctx context.Context, id int64) (states.State, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, ok := s.data[id]
	if !ok {
		return states.Unknown, fmt.Errorf("no such user")
	}

	return user.State, nil
}

func (s *Sessions) StartNewQuizForUser(ctx context.Context, id int64, newParty NewParty) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, ok := s.data[id]
	if !ok {
		return fmt.Errorf("no such user")
	}

	return user.StartNewQuiz(newParty)
}

func (s *Sessions) GetCurrentQuestionForUser(ctx context.Context, id int64) (bool, Question, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, ok := s.data[id]
	if !ok {
		return false, Question{}, fmt.Errorf("no such user")
	}

	questExists, question, err := user.GetCurrentQuestion()
	if err != nil {
		return false, Question{}, fmt.Errorf("can't get current question: %v", err)
	}
	return questExists, question, nil
}
