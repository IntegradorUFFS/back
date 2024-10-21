package materialServices

import (
	"back/internal/store/api/helper"
	materialTypes "back/internal/store/api/types/material"
	unitTypes "back/internal/store/api/types/unit"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func Autocomplete(w http.ResponseWriter, r *http.Request, p materialTypes.T_params, url_q materialTypes.T_autocompleteQuery, q *pgstore.Queries) {
	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	type _response struct {
		Data []materialTypes.T_autocompleteResponseBody `json:"data"`
	}
	parsed_uuid, err := uuid.Parse(url_q.ID.String())

	if err == nil && parsed_uuid != uuid.Nil {
		m, err := q.FindMaterialById(r.Context(), url_q.ID)

		if err != nil {
			data, _ := json.Marshal(_response{
				Data: []materialTypes.T_autocompleteResponseBody{},
			})

			w.Header().Set("Content-Type", "application/json")

			w.Write(data)
			return
		}

		unit, err := q.FindUnitById(r.Context(), m.UnitID)

		if err != nil {
			helper.HandleErrorMessage(w, err, "Unit")
			return
		}

		_materials := []materialTypes.T_autocompleteResponseBody{}
		_materials = append(_materials, materialTypes.T_autocompleteResponseBody{
			Name: m.Name,
			ID:   m.ID.String(),
			Unit: unitTypes.T_responseBody{
				ID:        unit.ID.String(),
				Name:      unit.Name,
				ShortName: unit.ShortName.String,
			},
		})

		data, _ := json.Marshal(_response{
			Data: _materials,
		})

		w.Header().Set("Content-Type", "application/json")

		w.Write(data)
		return
	}

	materials, err := q.AutocompleteMaterialByLikeName(r.Context(), url_q.Search)

	if err != nil {
		helper.HandleErrorMessage(w, err, "Location")
		return
	}

	_materials := []materialTypes.T_autocompleteResponseBody{}
	for _, m := range materials {
		unit, err := q.FindUnitById(r.Context(), m.UnitID)

		if err != nil {
			helper.HandleErrorMessage(w, err, "Unit")
			return
		}

		_materials = append(_materials, materialTypes.T_autocompleteResponseBody{
			Name: m.Name,
			ID:   m.ID.String(),
			Unit: unitTypes.T_responseBody{
				ID:        unit.ID.String(),
				Name:      unit.Name,
				ShortName: unit.ShortName.String,
			},
		})
	}

	data, _ := json.Marshal(_response{
		Data: _materials,
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
