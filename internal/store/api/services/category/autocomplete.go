package categoryServices

import (
	"back/internal/store/api/helper"
	categoryTypes "back/internal/store/api/types/category"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func Autocomplete(w http.ResponseWriter, r *http.Request, p categoryTypes.T_params, url_q categoryTypes.T_autocompleteQuery, q *pgstore.Queries) {
	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	type _response struct {
		Data []categoryTypes.T_responseBody `json:"data"`
	}
	parsed_uuid, err := uuid.Parse(url_q.ID.String())

	if err == nil && parsed_uuid != uuid.Nil {
		c, err := q.FindCategoryById(r.Context(), url_q.ID)

		if err != nil {
			data, _ := json.Marshal(_response{
				Data: []categoryTypes.T_responseBody{},
			})

			w.Header().Set("Content-Type", "application/json")

			w.Write(data)
			return
		}

		_categories := []categoryTypes.T_responseBody{}
		_categories = append(_categories, categoryTypes.T_responseBody{
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

	categories, err := q.AutocompleteCategoryByLikeName(r.Context(), url_q.Search)

	if err != nil {
		helper.HandleErrorMessage(w, err, "Category")
		return
	}

	_categories := []categoryTypes.T_responseBody{}
	for _, c := range categories {
		_categories = append(_categories, categoryTypes.T_responseBody{
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
