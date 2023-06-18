package domain

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain/states"
	"github.com/dimayasha7123/quiz_service/utils/logger"
	"strconv"
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
		data[user.ID] = NewUser(user.ID, user.QSID, user.Name)
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

func (s *Sessions) UserExistsByID(ctx context.Context, id int64) bool {
	s.mutex.RLock()
	_, ok := s.data[id]
	s.mutex.RUnlock()

	return ok
}

func (s *Sessions) UserByID(ctx context.Context, id int64) (User, error) {
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

	user := NewUser(id, qsid, name)

	s.data[user.TelegramID] = user

	return nil
}

func (s *Sessions) UserState(ctx context.Context, id int64) (states.State, error) {
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

func (s *Sessions) CurrentQuestionForUser(ctx context.Context, id int64) (bool, Question, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, ok := s.data[id]
	if !ok {
		return false, Question{}, fmt.Errorf("no such user")
	}

	questExists, question, err := user.CurrentQuestion()
	if err != nil {
		return false, Question{}, fmt.Errorf("can't get current question: %v", err)
	}
	return questExists, question, nil
}

func (s *Sessions) ConfirmQuestionForUser(ctx context.Context, id int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, ok := s.data[id]
	if !ok {
		return fmt.Errorf("no such user")
	}

	err := user.ConfirmQuestion()
	if err != nil {
		return fmt.Errorf("can't confirm question for user with id = %v: %v", id, err)
	}

	return nil
}

func (s *Sessions) EndQuizForUser(ctx context.Context, id int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, ok := s.data[id]
	if !ok {
		return fmt.Errorf("no such user")
	}

	err := user.EndQuiz()
	if err != nil {
		return fmt.Errorf("can't check if user can end quiz: %v", err)
	}

	return nil
}

func (s *Sessions) GetName(ctx context.Context, name string) string {
	id, err := strconv.ParseInt(name, 10, 64)
	if err != nil {
		logger.Log.Errorf("can't convert ID = %v from string to int: %v", name, err)
		return name
	}

	user, err := s.UserByID(ctx, id)
	if err != nil {
		logger.Log.Errorf("can't get user by id = %v from sessions: %v", id, err)
		return name
	}

	return user.Name
}

func (s *Sessions) SwitchAnswerForUser(ctx context.Context, id int64, ansID int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, ok := s.data[id]
	if !ok {
		return fmt.Errorf("no such user")
	}

	err := user.SwitchAnswer(ansID)
	if err != nil {
		return fmt.Errorf("can't switch answer with id = %v: %v", id, err)
	}

	return nil
}

func (s *Sessions) BreakSessionForUser(ctx context.Context, id int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, ok := s.data[id]
	if !ok {
		return fmt.Errorf("no such user")
	}

	err := user.BreakSession()
	if err != nil {
		return fmt.Errorf("can't break session for user with id = %v", err)
	}

	return nil
}
