package quizApi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dimayasha7123/quiz_service/internal/models"
	"io/ioutil"
)

func boolStrOrNilToBool(s string) bool {
	if s == "true" {
		return true
	}
	return false
}

func (qac *QuizPartyApiClient) GetParty(tag string) (*models.Party, error) {
	getUrl := fmt.Sprintf("https://quizapi.io/api/v1/questions?apiKey=%s&tags=%s&limit=%d", qac.apiKey, tag, questCount)
	resp, err := qac.cl.Get(getUrl)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	pd := PartyData{}
	err = json.Unmarshal(bytes, &pd)
	if err != nil {
		return nil, err
	}

	if len(pd) != questCount {
		return nil, errors.New(fmt.Sprintf("get %d questions instead of %d", len(pd), questCount))
	}

	questions := make([]models.Question, questCount)

	for i, q := range pd {

		answers := make([]models.Answer, 0, maxAnsCount)

		if q.Answers.AnswerA != "" {
			answers = append(answers, models.Answer{Title: q.Answers.AnswerA})
		}
		if q.Answers.AnswerB != "" {
			answers = append(answers, models.Answer{Title: q.Answers.AnswerB})
		}
		if q.Answers.AnswerC != "" {
			answers = append(answers, models.Answer{Title: q.Answers.AnswerC})
		}
		if q.Answers.AnswerD != "" {
			answers = append(answers, models.Answer{Title: q.Answers.AnswerD})
		}
		if q.Answers.AnswerE != "" {
			answers = append(answers, models.Answer{Title: q.Answers.AnswerE})
		}
		if q.Answers.AnswerF != "" {
			answers = append(answers, models.Answer{Title: q.Answers.AnswerF})
		}

		answersCorrect := make([]bool, maxAnsCount)

		if boolStrOrNilToBool(q.CorrectAnswers.AnswerACorrect) {
			answersCorrect[0] = true
		}
		if boolStrOrNilToBool(q.CorrectAnswers.AnswerBCorrect) {
			answersCorrect[1] = true
		}
		if boolStrOrNilToBool(q.CorrectAnswers.AnswerCCorrect) {
			answersCorrect[2] = true
		}
		if boolStrOrNilToBool(q.CorrectAnswers.AnswerDCorrect) {
			answersCorrect[3] = true
		}
		if boolStrOrNilToBool(q.CorrectAnswers.AnswerECorrect) {
			answersCorrect[4] = true
		}
		if boolStrOrNilToBool(q.CorrectAnswers.AnswerFCorrect) {
			answersCorrect[5] = true
		}

		for j := range answers {
			answers[j].Correct = answersCorrect[j]
		}

		tags := make([]string, len(q.Tags))

		for j, t := range q.Tags {
			tags[j] = t.Name
		}

		questions[i] = models.Question{
			Title:   q.Question,
			Tags:    tags,
			Answers: answers,
		}
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return &models.Party{Questions: questions}, nil
}
