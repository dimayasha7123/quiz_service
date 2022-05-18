package repository

import (
	"context"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/models"
)

func (r *repository) GetUserQuizTop(ctx context.Context, quizID, userID int64) (models.SingleTop, error) {
	userResultQuery := `
	select place, name, mp
		from (
		select ROW_NUMBER() over () as place, id, q.name, q.mp
			from (
			select ua.id id, ua.name name, max(points) mp
				from party
				join user_account ua on ua.id = party.user_account_id
				where completed
				and quiz_id = $1
				group by ua.id
				order by mp desc) as q) as qq
				where qq.id = $2;
	`

	singleTop := models.SingleTop{}
	err := r.pool.QueryRow(ctx, userResultQuery, quizID, userID).Scan(
		&singleTop.UserResults.Place,
		&singleTop.UserResults.Name,
		&singleTop.UserResults.Points,
	)
	if err != nil {
		return models.SingleTop{}, err
	}

	gTop, err := r.GetQuizTop(ctx, quizID)
	if err != nil {
		return models.SingleTop{}, err
	}

	singleTop.GlobalTop = gTop

	return singleTop, nil
}
