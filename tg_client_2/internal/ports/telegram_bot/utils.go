package telegram_bot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (s *service) sendMessage(ctx context.Context, chatID int64, text string, markup *replyMarkup) error {
	postURL := url.URL{
		Scheme: "https",
		Host:   "api.telegram.org",
		Path:   fmt.Sprintf("/bot%s/sendMessage", s.apiKey),
	}

	query := postURL.Query()
	query.Add("disable_web_page_preview", "1")
	query.Add("chat_id", fmt.Sprint(chatID))
	query.Add("text", text)
	if markup != nil {
		bytes, err := json.Marshal(markup)
		if err != nil {
			return fmt.Errorf("can't marshal reply markup: %v", err)
		}
		query.Add("reply_markup", string(bytes))
	}
	postURL.RawQuery = query.Encode()

	err := s.makePostReq(postURL.String())
	if err != nil {
		return fmt.Errorf("can't make post req: %v", err)
	}

	return nil
}

func (s *service) editMessage(ctx context.Context, chatID, messID int64, text string) error {
	postURL := url.URL{
		Scheme: "https",
		Host:   "api.telegram.org",
		Path:   fmt.Sprintf("/bot%s/editMessageText", s.apiKey),
	}

	query := postURL.Query()
	query.Add("disable_web_page_preview", "1")
	query.Add("chat_id", fmt.Sprint(chatID))
	query.Add("message_id", fmt.Sprint(messID))
	query.Add("text", text)
	postURL.RawQuery = query.Encode()

	err := s.makePostReq(postURL.String())
	if err != nil {
		return fmt.Errorf("can't make post req: %v", err)
	}

	return nil
}

func (s *service) editReplyMarkup(ctx context.Context, chatID, messID int64, markup replyMarkup) error {
	postURL := url.URL{
		Scheme: "https",
		Host:   "api.telegram.org",
		Path:   fmt.Sprintf("/bot%s/editMessageReplyMarkup", s.apiKey),
	}

	query := postURL.Query()
	query.Add("chat_id", fmt.Sprint(chatID))
	query.Add("message_id", fmt.Sprint(messID))

	bytes, err := json.Marshal(markup)
	if err != nil {
		return fmt.Errorf("can't marshal reply markup: %v", err)
	}
	query.Add("reply_markup", string(bytes))

	postURL.RawQuery = query.Encode()

	err = s.makePostReq(postURL.String())
	if err != nil {
		return fmt.Errorf("can't make post req: %v", err)
	}

	return nil
}

func (s *service) editMessageAndReplyMarkup(ctx context.Context, chatID, messID int64, text string, markup replyMarkup) error {
	postURL := url.URL{
		Scheme: "https",
		Host:   "api.telegram.org",
		Path:   fmt.Sprintf("/bot%s/editMessageText", s.apiKey),
	}

	query := postURL.Query()
	query.Add("disable_web_page_preview", "1")
	query.Add("chat_id", fmt.Sprint(chatID))
	query.Add("message_id", fmt.Sprint(messID))
	query.Add("text", text)
	bytes, err := json.Marshal(markup)
	if err != nil {
		return fmt.Errorf("can't marshal reply markup: %v", err)
	}
	query.Add("reply_markup", string(bytes))
	postURL.RawQuery = query.Encode()

	err = s.makePostReq(postURL.String())
	if err != nil {
		return fmt.Errorf("can't make post req: %v", err)
	}

	return nil
}

func (s *service) makePostReq(postURL string) error {
	resp, err := s.httpClient.Post(postURL, "text/plain", nil)
	if err != nil {
		return fmt.Errorf("can't make post req by url = %v: %v", postURL, err)
	}
	if resp == nil {
		return fmt.Errorf("nil resp")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("didn't get ok status, get %v status", resp.StatusCode)
	}

	return nil
}

func (s *service) sendCanon(ctx context.Context, chatID int64) error {
	canonURL := "https://media.tenor.com/GESYbde0oaYAAAAd/canon-event-lochan-bwefi.gif"
	postURL := url.URL{
		Scheme: "https",
		Host:   "api.telegram.org",
		Path:   fmt.Sprintf("/bot%s/sendAnimation", s.apiKey),
	}

	query := postURL.Query()
	query.Add("chat_id", fmt.Sprint(chatID))
	query.Add("animation", canonURL)
	postURL.RawQuery = query.Encode()

	err := s.makePostReq(postURL.String())
	if err != nil {
		return fmt.Errorf("can't make post req: %v", err)
	}

	return nil
}
