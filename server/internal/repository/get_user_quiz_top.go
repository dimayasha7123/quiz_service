package repository

import (
	"context"
	"github.com/dimayasha7123/quiz_service/server/internal/models"
)

func (r *repository) GetUserQuizTop(ctx context.Context, quizID, userID int64) (models.SingleTop, error) {

	countFinishedQuery := `
	select count(*) from party where user_account_id = $2 and quiz_id = $1 and completed;
	`

	var countFinished int64
	err := r.pool.QueryRow(ctx, countFinishedQuery, quizID, userID).Scan(&countFinished)
	if err != nil {
		return models.SingleTop{}, err
	}

	singleTop := models.SingleTop{}

	if countFinished != 0 {

		userResultQuery := `
	select place, name, mp
		from (
		select ROW_NUMBER() over () as place, id, q.name, q.mp
			from (
			select ua.id id, ua.name, max(points) mp
				from party
				join user_account ua on ua.id = party.user_account_id
				where completed
				and quiz_id = $1
				group by ua.id
				order by mp desc) as q) as qq
				where qq.id = $2;
	`
		err := r.pool.QueryRow(ctx, userResultQuery, quizID, userID).Scan(
			&singleTop.UserResults.Place,
			&singleTop.UserResults.Name,
			&singleTop.UserResults.Points,
		)
		if err != nil {
			return models.SingleTop{}, err
		}
	}

	gTop, err := r.GetQuizTop(ctx, quizID)
	if err != nil {
		return models.SingleTop{}, err
	}

	singleTop.GlobalTop = gTop

	return singleTop, nil
}
