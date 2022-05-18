package repository

import "context"

func (r *repository) GetRightAnswers(ctx context.Context, partyID int64) ([][]int32, error) {

	query := `
	select party_question.quest_order_number, a.order_number
	from party_question
         join question q on q.id = party_question.question_id
         join answer a on q.id = a.question_id
	where party_question.party_id = $1 and a.correct
	order by party_question.quest_order_number, a.order_number;
	`
	rows, err := r.pool.Query(ctx, query, partyID)
	if err != nil {
		return nil, err
	}

	rightAnswers := make([][]int32, 10) // TODO вынести количество вопросов в конфиг

	for rows.Next() {
		var q, a int32
		err = rows.Scan(&q, &a)
		if err != nil {
			return nil, err
		}
		rightAnswers[q] = append(rightAnswers[q], a)
	}

	return rightAnswers, nil
}
