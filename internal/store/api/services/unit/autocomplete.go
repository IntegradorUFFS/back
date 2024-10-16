package unitServices

import (
	"back/internal/store/api/helper"
	unitTypes "back/internal/store/api/types/unit"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func Autocomplete(w http.ResponseWriter, r *http.Request, p unitTypes.T_params, url_q unitTypes.T_autocompleteQuery, q *pgstore.Queries) {
	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	type _response struct {
		Data []unitTypes.T_responseBody `json:"data"`
	}
	parsed_uuid, err := uuid.Parse(url_q.ID.String())

	if err == nil && parsed_uuid != uuid.Nil {
		c, err := q.FindUnitById(r.Context(), url_q.ID)

		if err != nil {
			data, _ := json.Marshal(_response{
				Data: []unitTypes.T_responseBody{},
			})

			w.Header().Set("Content-Type", "application/json")

			w.Write(data)
			return
		}

		_units := []unitTypes.T_responseBody{}
		_units = append(_units, unitTypes.T_responseBody{
			ShortName: c.ShortName.String,
			Name:      c.Name,
			ID:        c.ID.String(),
		})

		data, _ := json.Marshal(_response{
			Data: _units,
		})

		w.Header().Set("Content-Type", "application/json")

		w.Write(data)
		return
	}

	units, err := q.AutocompleteUnitByLikeName(r.Context(), url_q.Search)

	if err != nil {
		helper.HandleErrorMessage(w, err, "Cunit")
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

	data, _ := json.Marshal(_response{
		Data: _units,
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
