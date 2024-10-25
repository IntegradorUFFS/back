package transactionServices

import (
	"back/internal/store/api/helper"
	transactionTypes "back/internal/store/api/types/transaction"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

func Read(w http.ResponseWriter, r *http.Request, p transactionTypes.T_params, url_q transactionTypes.T_readQuery, q *pgstore.Queries) {
	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	offset := url_q.Page * url_q.PerPage
	limit := url_q.PerPage

	filter := ""

	type _response struct {
		Data []transactionTypes.T_responseCleanBody `json:"data"`
		Meta transactionTypes.T_responseMeta        `json:"meta"`
	}

	filter_count := 1

	var filters_arr []any

	if url_q.FilterMaterialID != "" {
		if filter_count == 1 {
			filter += " WHERE"
		} else {
			filter += " AND"
		}
		filter += " material_id = $" + fmt.Sprint(filter_count)
		filter_count += 1

		filters_arr = append(filters_arr, url_q.FilterMaterialID)
	}
	
	if url_q.FilterOriginLocationID != "" {
		if filter_count == 1 {
			filter += " WHERE"
		} else {
			filter += " AND"
		}
		filter += " origin_id = $" + fmt.Sprint(filter_count)
		filter_count += 1

		filters_arr = append(filters_arr, url_q.FilterOriginLocationID)
	}

	if url_q.FilterDestinyLocationID != "" {
		if filter_count == 1 {
			filter += " WHERE"
		} else {
			filter += " AND"
		}
		filter += " destiny_id = $" + fmt.Sprint(filter_count)
		filter_count += 1

		filters_arr = append(filters_arr, url_q.FilterDestinyLocationID)
	}

	if url_q.FilterType != "" {
		if filter_count == 1 {
			filter += " WHERE"
		} else {
			filter += " AND"
		}
		filter += " type = $" + fmt.Sprint(filter_count)
		filter_count += 1

		filters_arr = append(filters_arr, url_q.FilterType)
	}

	size_filters := filters_arr

	filters_arr = append(filters_arr, limit)
	filters_arr = append(filters_arr, offset)

	transactions, err := q.C_FetchPaginatedTransactionsWithJson(r.Context(), `SELECT json_build_object(
    'id', transaction.id,
    'quantity', transaction.quantity,
    'created_at', transaction.created_at,
    'type', transaction.type,
    'material', json_build_object(
        'id', material.id,
        'name', material.name,
        'description', material.description,
        'quantity', material.quantity
    ),
    'origin', json_build_object(
        'id', origin.id,
        'name', origin.name
    ),
    'destiny', json_build_object(
        'id', destiny.id,
        'name', destiny.name
    )
)
FROM transaction
LEFT JOIN material ON transaction.material_id = material.id
LEFT JOIN location origin ON transaction.origin_location_id = origin.id
LEFT JOIN location destiny ON transaction.destiny_location_id = destiny.id
`+filter+
		" ORDER BY "+url_q.SortColumn+" "+url_q.SortDirection+", created_at DESC LIMIT $" + fmt.Sprint(filter_count) + " OFFSET $" + fmt.Sprint(filter_count + 1), filters_arr)

	if err != nil {
		helper.HandleErrorMessage(w, err, "None transaction")
		return
	}

	size, err := q.C_GetTableSize(r.Context(), `SELECT count(*) AS exact_count FROM transaction`+filter, size_filters)

	if err != nil {
		helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
		return
	}

	total_pages := math.Ceil(float64(size) / float64(limit))

	var _transactions []transactionTypes.T_responseCleanBody

	for _, t := range transactions {
		var temp transactionTypes.T_responseCleanBody

		json.Unmarshal(t, &temp)
		_transactions = append(_transactions, temp)
	}

	data, _ := json.Marshal(_response{
		Data: _transactions,
		Meta: transactionTypes.T_responseMeta{
			Page:       url_q.Page,
			PerPage:    url_q.PerPage,
			TotalPages: int32(total_pages),
			Total:      int32(size),
		},
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
