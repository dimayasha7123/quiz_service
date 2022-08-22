package repository

import (
	"context"
	"github.com/dimayasha7123/quiz_service/internal/models"
)

func (r *repository) GetQuizTop(ctx context.Context, quizID int64) (models.GlobalTop, error) {
	globalTopQuery := `
	select ua.name, max(points) mp
	from party
			join user_account ua on ua.id = party.user_account_id
	where completed
	  and quiz_id = $1
	group by ua.id
	order by mp desc
	limit 3;
	`
	rows, err := r.pool.Query(ctx, globalTopQuery, quizID)
	if err != nil {
		return models.GlobalTop{}, err
	}

	place := int64(1)
	gTop := models.GlobalTop{}

	for rows.Next() {
		var name string
		var pts int32
		err = rows.Scan(&name, &pts)
		gTop.Results = append(gTop.Results, models.PartyResults{
			Name:   name,
			Points: pts,
			Place:  place,
		})
		place++
	}

	return gTop, nil
}
