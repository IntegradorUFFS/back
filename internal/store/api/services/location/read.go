package locationServices

import (
	"back/internal/store/api/helper"
	locationTypes "back/internal/store/api/types/location"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
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

	type _response struct {
		Data []locationTypes.T_responseBody `json:"data"`
		Meta locationTypes.T_responseMeta   `json:"meta"`
	}

	categories, err := q.FetchPaginatedLocations(r.Context(), pgstore.FetchPaginatedLocationsParams{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		helper.HandleErrorMessage(w, err, "None location")
		return
	}

	size, err := q.GetLocationTableSize(r.Context())

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
