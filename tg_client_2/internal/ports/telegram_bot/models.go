package telegram_bot

type updates struct {
	Ok     bool     `json:"ok"`
	Result []update `json:"result"`
}

type update struct {
	UpdateID      int64          `json:"update_id"`
	Message       *message       `json:"message"`
	CallbackQuery *callbackQuery `json:"callback_query"`
}

type message struct {
	MessageID int64 `json:"message_id"`
	From      struct {
		ID           int64  `json:"id"`
		IsBot        bool   `json:"is_bot"`
		FirstName    string `json:"first_name"`
		Username     string `json:"username"`
		LanguageCode string `json:"language_code"`
	} `json:"from"`
	Chat struct {
		ID        int64  `json:"id"`
		FirstName string `json:"first_name"`
		Username  string `json:"username"`
		Type      string `json:"type"`
	} `json:"chat"`
	Date int64  `json:"date"`
	Text string `json:"text"`
}

type callbackQuery struct {
	ID   string `json:"id"`
	From struct {
		ID           int64  `json:"id"`
		IsBot        bool   `json:"is_bot"`
		FirstName    string `json:"first_name"`
		Username     string `json:"username"`
		LanguageCode string `json:"language_code"`
	} `json:"from"`
	Message struct {
		MessageID int64 `json:"message_id"`
		From      struct {
			ID        int64  `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID        int64  `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date        int64  `json:"date"`
		Text        string `json:"text"`
		ReplyMarkup struct {
			InlineKeyboard [][]struct {
				Text         string `json:"text"`
				CallbackData string `json:"callback_data"`
			} `json:"inline_keyboard"`
		} `json:"reply_markup"`
	} `json:"message"`
	ChatInstance string `json:"chat_instance"`
	Data         string `json:"data"`
}

type replyMarkup struct {
	InlineKeyboard [][]replyButton `json:"inline_keyboard"`
}

type replyButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}
