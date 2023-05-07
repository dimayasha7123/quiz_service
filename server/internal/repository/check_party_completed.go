package repository

import "context"

func (r *repository) CheckPartyCompleted(ctx context.Context, partyID int64) (bool, error) {
	query := `
	select completed from party where id = $1;
	`
	var comp bool
	err := r.pool.QueryRow(ctx, query, partyID).Scan(&comp)
	if err != nil {
		return false, err
	}
	return comp, nil
}
