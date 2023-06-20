package telegram_bot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dimayasha7123/quiz_service/utils/logger"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	errorText = "Oops, get some error ("
)

func (s *service) Run(ctx context.Context) error {
	lastUpdateID := int64(0)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			newLastUpdateID, err := s.runIteration(ctx, lastUpdateID)
			if err != nil {
				logger.Log.Errorf("Get error while running: %v", err)
			}
			lastUpdateID = newLastUpdateID
		}

		time.Sleep(s.delay)
	}
}

func (s *service) runIteration(ctx context.Context, lastUpdateID int64) (int64, error) {
	reqURL := url.URL{
		Scheme: "https",
		Host:   "api.telegram.org",
		Path:   fmt.Sprintf("/bot%s/getUpdates", s.apiKey),
	}
	query := reqURL.Query()
	query.Add("offset", fmt.Sprint(lastUpdateID+1))
	reqURL.RawQuery = query.Encode()

	resp, err := s.httpClient.Get(reqURL.String())
	if err != nil {
		return lastUpdateID, fmt.Errorf("can't make get req by url = %v: %v", reqURL, err)
	}
	if resp == nil {
		return lastUpdateID, fmt.Errorf("nil resp")
	}
	if resp.StatusCode != http.StatusOK {
		return lastUpdateID, fmt.Errorf("didn't get ok status, get %v status", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return lastUpdateID, fmt.Errorf("can't read resp body: %v", err)
	}

	err = resp.Body.Close()
	if err != nil {
		logger.Log.Errorf("can't close resp body, but continue: %v", err)
	}

	var ups updates
	err = json.Unmarshal(bytes, &ups)
	if err != nil {
		return lastUpdateID, fmt.Errorf("can't unmarshall resp body bytes: %v", err)
	}

	if !ups.Ok {
		return lastUpdateID, fmt.Errorf("didn't get ok in resp body")
	}

	// TODO: обработать кастомных ошибок и их проверка для ответа пользователю, в чем проблема (no such command)

	for _, up := range ups.Result {
		up := up

		go func(ctx context.Context) {
			var err error
			var errID int64

			switch {
			case up.CallbackQuery != nil:
				err = s.routeCallbackQuery(ctx, *up.CallbackQuery)
				errID = up.CallbackQuery.From.ID
			case up.Message != nil:
				err = s.routeMessage(ctx, *up.Message)
				errID = up.Message.From.ID
			default:
				logger.Log.Errorf("no new message or callback query in update with id = %")
			}

			if err != nil {
				logger.Log.Errorf("can't route request: %v", err)
				err = s.sendMessage(ctx, errID, errorText, nil)
				if err != nil {
					logger.Log.Errorf("can't send error message to user with id = %v: %v", errID, err)
				}
			}
		}(ctx)
	}

	if len(ups.Result) != 0 {
		lastUpdateID = ups.Result[len(ups.Result)-1].UpdateID
	}

	return lastUpdateID, nil
}

// TODO: добавить логгирование запросов и ответов

// 	команды:
//		/start		приветствие и "чтобы быстро начать воспользуетесь /quizzes"
//		/help		подробное описание и опять /quizzes
//		/quizzes 	список кнопок (1) с названиями квизов
//		/break		прерывает текущий квиз, если есть
//
// 	кнопки:
//		(1) квиз - возвращает кнопки (2) "принять участие", (3) "посмотреть топ", (4) "назад"
//		(2) принять участие - запускает квиз, возвращает вопрос и кнопки (5) с ответами и (6) "подтвердить"
//		(3) показывает топ по квизу и кнопку (7) "назад" (к квизу)
//		(4) -> /quizzes
//		(5) меняет выбранный ответ, показывает вопрос (2)
//		(6) показывает следующий вопрос или завершает квиз и показывает рейтинг и кнопку (к квизам), т.е -> /quizzes
//
//
//
//
