package pgstore

import "context"

func (q *Queries) C_UpdateUser(ctx context.Context, query string, arg []any) (UpdateUserRow, error) {
	row := q.db.QueryRow(ctx, query, arg...)
	var i UpdateUserRow
	err := row.Scan(
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.Role,
	)
	return i, err
}

func (q *Queries) C_UpdateUnit(ctx context.Context, query string, arg []any) (Unit, error) {
	row := q.db.QueryRow(ctx, query, arg...)
	var i Unit
	err := row.Scan(
		&i.Name,
		&i.ShortName,
	)
	return i, err
}

func (q *Queries) C_FetchPaginatedCategories(ctx context.Context, query string, arg []any) ([]Category, error) {
	rows, err := q.db.Query(ctx, query, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Category
	for rows.Next() {
		var i Category
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *Queries) C_GetTableSize(ctx context.Context, query string, arg []any) (int64, error) {
	row := q.db.QueryRow(ctx, query, arg...)
	var exact_count int64
	err := row.Scan(&exact_count)
	return exact_count, err
}


func (q *Queries) C_FetchPaginatedLocations(ctx context.Context, query string, arg []any) ([]Location, error) {
	rows, err := q.db.Query(ctx, query, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Location
	for rows.Next() {
		var i Location
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *Queries) C_FetchPaginatedUnits(ctx context.Context, query string, arg []any) ([]Unit, error) {
	rows, err := q.db.Query(ctx, query, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Unit
	for rows.Next() {
		var i Unit
		if err := rows.Scan(&i.ID, &i.Name, &i.ShortName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *Queries) C_FetchPaginatedUsers(ctx context.Context, query string, arg []any) ([]FetchPaginatedUsersRow, error) {
	rows, err := q.db.Query(ctx, query, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FetchPaginatedUsersRow
	for rows.Next() {
		var i FetchPaginatedUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.FirstName,
			&i.LastName,
			&i.Role,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *Queries) C_FetchPaginatedMaterials(ctx context.Context, query string, arg []any) ([]Material, error) {
	rows, err := q.db.Query(ctx, query, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Material
	for rows.Next() {
		var i Material
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Quantity,
			&i.CategoryID,
			&i.UnitID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *Queries) C_UpdateMaterial(ctx context.Context, query string, arg []any) (UpdateMaterialRow, error) {
	row := q.db.QueryRow(ctx, query, arg...)
	var i UpdateMaterialRow
	err := row.Scan(
		&i.Name,
		&i.Description,
		&i.Quantity,
		&i.CategoryID,
		&i.UnitID,
	)
	return i, err
}

func (q *Queries) C_FetchPaginatedLocationMaterials(ctx context.Context, query string, arg []any) ([]LocationMaterial, error) {
	rows, err := q.db.Query(ctx, query, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []LocationMaterial
	for rows.Next() {
		var i LocationMaterial
		if err := rows.Scan(
			&i.ID,
			&i.Quantity,
			&i.MaterialID,
			&i.LocationID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *Queries) C_CreateTransaction(ctx context.Context, arg CreateTransactionParams) (CreateTransactionRow, error) {
	row := q.db.QueryRow(ctx, createTransaction,
		arg.Quantity,
		arg.Type,
		arg.DestinyLocationID,
		arg.OriginLocationID,
		arg.MaterialID,
	)
	var i CreateTransactionRow
	err := row.Scan(
		&i.ID,
		&i.Quantity,
		&i.Type,
		&i.OriginLocationID,
		&i.DestinyLocationID,
		&i.MaterialID,
		&i.CreatedAt,
	)
	return i, err
}

func (q *Queries) C_CreateTransactionWithDL(ctx context.Context, arg CreateTransactionWithDLParams) (CreateTransactionRow, error) {
	row := q.db.QueryRow(ctx, createTransactionWithDL,
		arg.Quantity,
		arg.Type,
		arg.DestinyLocationID,
		arg.MaterialID,
	)
	var i CreateTransactionRow
	err := row.Scan(
		&i.ID,
		&i.Quantity,
		&i.Type,
		&i.OriginLocationID,
		&i.DestinyLocationID,
		&i.MaterialID,
		&i.CreatedAt,
	)
	return i, err
}

func (q *Queries) C_CreateTransactionWithOL(ctx context.Context, arg CreateTransactionWithOLParams) (CreateTransactionRow, error) {
	row := q.db.QueryRow(ctx, createTransactionWithOL,
		arg.Quantity,
		arg.Type,
		arg.OriginLocationID,
		arg.MaterialID,
	)
	var i CreateTransactionRow
	err := row.Scan(
		&i.ID,
		&i.Quantity,
		&i.Type,
		&i.OriginLocationID,
		&i.DestinyLocationID,
		&i.MaterialID,
		&i.CreatedAt,
	)
	return i, err
}

func (q *Queries) C_FetchPaginatedTransactions(ctx context.Context, query string, arg FetchPaginatedTransactionsParams) ([]FetchPaginatedTransactionsRow, error) {
	rows, err := q.db.Query(ctx, query, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FetchPaginatedTransactionsRow
	for rows.Next() {
		var i FetchPaginatedTransactionsRow
		if err := rows.Scan(
			&i.ID,
			&i.Quantity,
			&i.Type,
			&i.OriginLocationID,
			&i.DestinyLocationID,
			&i.MaterialID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *Queries) C_FetchPaginatedTransactionsWithJson(ctx context.Context, query string, arg []any) ([][]byte, error) {
	rows, err := q.db.Query(ctx, query, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items [][]byte
	for rows.Next() {
		var json_build_object []byte
		if err := rows.Scan(&json_build_object); err != nil {
			return nil, err
		}
		items = append(items, json_build_object)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
