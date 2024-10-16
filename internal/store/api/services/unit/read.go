package unitServices

import (
	"back/internal/store/api/helper"
	unitTypes "back/internal/store/api/types/unit"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"math"
	"net/http"
)

func Read(w http.ResponseWriter, r *http.Request, p unitTypes.T_params, url_q unitTypes.T_readQuery, q *pgstore.Queries) {
	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	offset := url_q.Page * url_q.PerPage
	limit := url_q.PerPage

	type _response struct {
		Data []unitTypes.T_responseBody `json:"data"`
		Meta unitTypes.T_responseMeta   `json:"meta"`
	}

	units, err := q.FetchPaginatedUnits(r.Context(), pgstore.FetchPaginatedUnitsParams{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		helper.HandleErrorMessage(w, err, "None unit")
		return
	}

	size, err := q.GetUnitTableSize(r.Context())

	if err != nil {
		helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
		return
	}

	_units := []unitTypes.T_responseBody{}
	for _, c := range units {
		_units = append(_units, unitTypes.T_responseBody{
			ShortName: c.ShortName.String,
			Name:      c.Name,
			ID:        c.ID.String(),
		})
	}

	total_pages := math.Ceil(float64(size) / float64(limit))

	data, _ := json.Marshal(_response{
		Data: _units,
		Meta: unitTypes.T_responseMeta{
			Page:       url_q.Page,
			PerPage:    url_q.PerPage,
			TotalPages: int32(total_pages),
			Total:      int32(size),
		},
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
