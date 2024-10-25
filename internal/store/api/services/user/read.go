package userServices

import (
	"back/internal/store/api/helper"
	userTypes "back/internal/store/api/types/user"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"
)

func Read(w http.ResponseWriter, r *http.Request, p userTypes.T_params, url_q userTypes.T_readQuery, q *pgstore.Queries) {
	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	offset := url_q.Page * url_q.PerPage
	limit := url_q.PerPage
	filter := ""

	type _response struct {
		Data []userTypes.T_responseBodyWithCreatedAt `json:"data"`
		Meta userTypes.T_responseMeta                `json:"meta"`
	}


	filter_count := 1

	var filters_arr []any

	if url_q.FilterFirstName != "" {
		if filter_count == 1 {
			filter += " WHERE"
		} else {
			filter += " AND"
		}
		filter += " first_name ~* $" + fmt.Sprint(filter_count)
		filter_count += 1

		filters_arr = append(filters_arr, url_q.FilterFirstName)
	}

	if url_q.FilterLastName != "" {
		if filter_count == 1 {
			filter += " WHERE"
		} else {
			filter += " AND"
		}
		filter += " last_name ~* $" + fmt.Sprint(filter_count)
		filter_count += 1

		filters_arr = append(filters_arr, url_q.FilterLastName)
	}

	if url_q.FilterEmail != "" {
		if filter_count == 1 {
			filter += " WHERE"
		} else {
			filter += " AND"
		}
		filter += " email ~* $" + fmt.Sprint(filter_count)
		filter_count += 1

		filters_arr = append(filters_arr, url_q.FilterEmail)
	}

	if p.RequesterRole == "manager" {
		if filter_count == 1 {
			filter += " WHERE"
		} else {
			filter += " AND"
		}
		filter += " role = 'viewer'"
	} else if url_q.FilterRole != "" {
		if filter_count == 1 {
			filter += " WHERE"
		} else {
			filter += " AND"
		}
		filter += " role = '" + url_q.FilterRole + "'"
	}

	size_filters := filters_arr

	filters_arr = append(filters_arr, limit)
	filters_arr = append(filters_arr, offset)

	_users := []userTypes.T_responseBodyWithCreatedAt{}

	users, err := q.C_FetchPaginatedUsers(r.Context(), "SELECT id, email, first_name, last_name, role, created_at FROM users"+filter+
		" ORDER BY "+url_q.SortColumn+" "+url_q.SortDirection+" LIMIT $" + fmt.Sprint(filter_count) + " OFFSET $" + fmt.Sprint(filter_count + 1), filters_arr)

	if err != nil {
		helper.HandleErrorMessage(w, err, "None user")
		return
	}

	size, err := q.C_GetTableSize(r.Context(), `SELECT count(*) AS exact_count FROM users`+filter,size_filters)

	if err != nil {
		helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
		return
	}

	for _, u := range users {
		_users = append(_users, userTypes.T_responseBodyWithCreatedAt{
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName.String,
			Role:      string(u.Role),
			ID:        u.ID.String(),
			CreatedAt: u.CreatedAt.Time.Format(time.RFC3339),
		})
	}

	total_pages := math.Ceil(float64(size) / float64(limit))

	data, _ := json.Marshal(_response{
		Data: _users,
		Meta: userTypes.T_responseMeta{
			Page:       url_q.Page,
			PerPage:    url_q.PerPage,
			TotalPages: int32(total_pages),
			Total:      int32(size),
		},
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
