package repository

import (
	"context"
	"github.com/dimayasha7123/quiz_service/server/internal/models"
)

func (r *repository) CompleteParty(ctx context.Context, partyID int64, points int32) (models.SingleTop, error) {

	updatePartyQuery := `
	update party
	set completed = true,
    	points    = $2
	where id = $1 returning id;
	`

	var tmp int64
	err := r.pool.QueryRow(ctx, updatePartyQuery, partyID, points).Scan(&tmp)
	if err != nil {
		return models.SingleTop{}, err
	}

	partyResultQuery := `
	select place, name, mp
	from (
			 select ROW_NUMBER() over () as place, id, qqq.name, qqq.mp
			 from (
					  select qq.id, ua3.name, qq.mp
					  from (select q.id, max(q.points) mp
							from (select ua.id, points
								  from party p
										   join user_account ua on ua.id = p.user_account_id
								  where completed
									and quiz_id = (select quiz_id from party where id = $1)
									and p.user_account_id != (select user_account_id from party where id = $1)
								  union
								  select ua2.id, points
								  from party p2
										   join user_account ua2 on ua2.id = p2.user_account_id
								  where p2.id = $1) as q
							group by q.id
							order by mp desc
							limit 3) as qq
							   join user_account ua3 on qq.id = ua3.id) as qqq) as qqqq
	where qqqq.id = (select user_account_id from party where id = $1);
	`
	singleTop := models.SingleTop{}
	err = r.pool.QueryRow(ctx, partyResultQuery, partyID).Scan(
		&singleTop.UserResults.Place,
		&singleTop.UserResults.Name,
		&singleTop.UserResults.Points,
	)
	if err != nil {
		return models.SingleTop{}, err
	}

	globalTopQuery := `
	select ua3.name, qq.mp
		from (select q.id, max(q.points) mp
			from (select ua.id, points
				from party p
				join user_account ua on ua.id = p.user_account_id
				where completed
				and quiz_id = (select quiz_id from party where id = $1)
					and p.user_account_id != (select user_account_id from party where id = $1)
						union
						select ua2.id, points
							from party p2
							join user_account ua2 on ua2.id = p2.user_account_id
							where p2.id = $1) as q
							group by q.id
							order by mp desc
							limit 3) as qq
							join user_account ua3 on qq.id = ua3.id;
	`

	rows, err := r.pool.Query(ctx, globalTopQuery, partyID)
	if err != nil {
		return models.SingleTop{}, err
	}

	place := int64(1)

	for rows.Next() {
		var name string
		var pts int32
		err = rows.Scan(&name, &pts)
		singleTop.GlobalTop.Results = append(singleTop.GlobalTop.Results, models.PartyResults{
			Name:   name,
			Points: pts,
			Place:  place,
		})
		place++
	}

	return singleTop, nil
}
