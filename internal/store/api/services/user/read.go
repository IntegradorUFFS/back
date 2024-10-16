package userServices

import (
	"back/internal/store/api/helper"
	userTypes "back/internal/store/api/types/user"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"math"
	"net/http"
)

func Read(w http.ResponseWriter, r *http.Request, p userTypes.T_params, url_q userTypes.T_readQuery, q *pgstore.Queries) {
	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	offset := url_q.Page * url_q.PerPage
	limit := url_q.PerPage

	type _response struct {
		Data []userTypes.T_responseBody `json:"data"`
		Meta userTypes.T_responseMeta   `json:"meta"`
	}

	_users := []userTypes.T_responseBody{}

	if p.RequesterRole == "manager" {
		users, err := q.FetchPaginatedUsersByRole(r.Context(), pgstore.FetchPaginatedUsersByRoleParams{
			Role:   "viewer",
			Offset: offset,
			Limit:  limit,
		})

		if err != nil {
			helper.HandleErrorMessage(w, err, "None user")
			return
		}

		size, err := q.GetRoledUserTableSize(r.Context(), "viewer")

		if err != nil {
			helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
			return
		}

		for _, u := range users {
			_users = append(_users, userTypes.T_responseBody{
				Email:     u.Email,
				FirstName: u.FirstName,
				LastName:  u.LastName.String,
				Role:      string(u.Role),
				ID:        u.ID.String(),
			})
		}

		total_pages := math.Ceil(float64(size) / float64(limit))

		data, _ := json.Marshal(_response{
			Data: []userTypes.T_responseBody(_users),
			Meta: userTypes.T_responseMeta{
				Page:       url_q.Page,
				PerPage:    url_q.PerPage,
				TotalPages: int32(total_pages),
				Total:      int32(size),
			},
		})

		w.Header().Set("Content-Type", "application/json")

		w.Write(data)
		return
	}

	users, err := q.FetchPaginatedUsers(r.Context(), pgstore.FetchPaginatedUsersParams{
		Offset: offset,
		Limit:  limit,
	})

	if err != nil {
		helper.HandleErrorMessage(w, err, "None user")
		return
	}

	size, err := q.GetUserTableSize(r.Context())

	if err != nil {
		helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
		return
	}

	for _, u := range users {
		_users = append(_users, userTypes.T_responseBody{
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName.String,
			Role:      string(u.Role),
			ID:        u.ID.String(),
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
