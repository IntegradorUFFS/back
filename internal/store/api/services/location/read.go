package locationServices

import (
	"back/internal/store/api/helper"
	locationTypes "back/internal/store/api/types/location"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

func Read(w http.ResponseWriter, r *http.Request, p locationTypes.T_params, url_q locationTypes.T_readQuery, q *pgstore.Queries) {
	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	offset := url_q.Page * url_q.PerPage
	limit := url_q.PerPage
	filter := ""

	type _response struct {
		Data []locationTypes.T_responseBody `json:"data"`
		Meta locationTypes.T_responseMeta   `json:"meta"`
	}


	filter_count := 1

	var filters_arr []any

	if url_q.FilterName != "" {
		if filter_count == 1 {
			filter += " WHERE"
		} else {
			filter += " AND"
		}
		filter += " name ~* $" + fmt.Sprint(filter_count)
		filter_count += 1

		filters_arr = append(filters_arr, url_q.FilterName)
	}

	size_filters := filters_arr

	filters_arr = append(filters_arr, limit)
	filters_arr = append(filters_arr, offset)


	categories, err := q.C_FetchPaginatedLocations(r.Context(), "SELECT id, name FROM location"+filter+
		" ORDER BY "+url_q.SortColumn+" "+url_q.SortDirection+" LIMIT $" + fmt.Sprint(filter_count) + " OFFSET $" + fmt.Sprint(filter_count + 1), filters_arr)

	if err != nil {
		helper.HandleErrorMessage(w, err, "None location")
		return
	}

	size, err := q.C_GetTableSize(r.Context(), `SELECT count(*) AS exact_count FROM location`+filter, size_filters)

	if err != nil {
		helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
		return
	}

	_categories := []locationTypes.T_responseBody{}
	for _, c := range categories {
		_categories = append(_categories, locationTypes.T_responseBody{
			Name: c.Name,
			ID:   c.ID.String(),
		})
	}

	total_pages := math.Ceil(float64(size) / float64(limit))

	data, _ := json.Marshal(_response{
		Data: _categories,
		Meta: locationTypes.T_responseMeta{
			Page:       url_q.Page,
			PerPage:    url_q.PerPage,
			TotalPages: int32(total_pages),
			Total:      int32(size),
		},
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
