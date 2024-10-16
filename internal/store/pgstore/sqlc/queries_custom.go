package pgstore

import "context"

func (q *Queries) C_UpdateUser(ctx context.Context, query string, arg UpdateUserParams) (UpdateUserRow, error) {
	row := q.db.QueryRow(ctx, query, arg.ID)
	var i UpdateUserRow
	err := row.Scan(
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.Role,
	)
	return i, err
}

func (q *Queries) C_UpdateUnit(ctx context.Context, query string, arg UpdateUnitParams) (Unit, error) {
	row := q.db.QueryRow(ctx, query, arg.ID)
	var i Unit
	err := row.Scan(
		&i.Name,
		&i.ShortName,
	)
	return i, err
}
