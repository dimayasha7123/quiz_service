package quizApi

import (
	"net/http"
	"time"
)

type PartyData []struct {
	ID          int    `json:"id"`
	Question    string `json:"question"`
	Description string `json:"description"`
	Answers     struct {
		AnswerA string `json:"answer_a"`
		AnswerB string `json:"answer_b"`
		AnswerC string `json:"answer_c"`
		AnswerD string `json:"answer_d"`
		AnswerE string `json:"answer_e"`
		AnswerF string `json:"answer_f"`
	} `json:"answers"`
	MultipleCorrectAnswers string `json:"multiple_correct_answers"`
	CorrectAnswers         struct {
		AnswerACorrect string `json:"answer_a_correct"`
		AnswerBCorrect string `json:"answer_b_correct"`
		AnswerCCorrect string `json:"answer_c_correct"`
		AnswerDCorrect string `json:"answer_d_correct"`
		AnswerECorrect string `json:"answer_e_correct"`
		AnswerFCorrect string `json:"answer_f_correct"`
	} `json:"correct_answers"`
	CorrectAnswer string `json:"correct_answer"`
	Explanation   string `json:"explanation"`
	Tip           string `json:"tip"`
	Tags          []struct {
		Name string `json:"name"`
	} `json:"tags"`
	Category   string `json:"category"`
	Difficulty string `json:"difficulty"`
}

const (
	ansCount   = 6
	questCount = 10
)

type quizPartyApiClient struct {
	apiKey string
	cl     http.Client
}

func New(apiKey string) *quizPartyApiClient {
	return &quizPartyApiClient{
		cl:     http.Client{Timeout: 10 * time.Second},
		apiKey: apiKey,
	}
}
