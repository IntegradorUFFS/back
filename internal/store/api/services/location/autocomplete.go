package locationServices

import (
	"back/internal/store/api/helper"
	locationTypes "back/internal/store/api/types/location"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func Autocomplete(w http.ResponseWriter, r *http.Request, p locationTypes.T_params, url_q locationTypes.T_autocompleteQuery, q *pgstore.Queries) {
	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	type _response struct {
		Data []locationTypes.T_responseBody `json:"data"`
	}
	parsed_uuid, err := uuid.Parse(url_q.ID.String())

	if err == nil && parsed_uuid != uuid.Nil {
		c, err := q.FindLocationById(r.Context(), url_q.ID)

		if err != nil {
			data, _ := json.Marshal(_response{
				Data: []locationTypes.T_responseBody{},
			})

			w.Header().Set("Content-Type", "application/json")

			w.Write(data)
			return
		}

		_categories := []locationTypes.T_responseBody{}
		_categories = append(_categories, locationTypes.T_responseBody{
			Name: c.Name,
			ID:   c.ID.String(),
		})

		data, _ := json.Marshal(_response{
			Data: _categories,
		})

		w.Header().Set("Content-Type", "application/json")

		w.Write(data)
		return
	}

	var categories []pgstore.Location

	if url_q.FilterMaterialID != uuid.Nil {
		categories, err = q.AutocompleteLocationByLikeNameWithMaterial(r.Context(), pgstore.AutocompleteLocationByLikeNameWithMaterialParams{
			MaterialID: url_q.FilterMaterialID,
			Name:       url_q.Search,
		})
	} else {
		categories, err = q.AutocompleteLocationByLikeName(r.Context(), url_q.Search)
	}
	if err != nil {
		helper.HandleErrorMessage(w, err, "Location")
		return
	}

	_categories := []locationTypes.T_responseBody{}
	for _, c := range categories {
		_categories = append(_categories, locationTypes.T_responseBody{
			Name: c.Name,
			ID:   c.ID.String(),
		})
	}

	data, _ := json.Marshal(_response{
		Data: _categories,
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
