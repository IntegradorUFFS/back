package materialServices

import (
	"back/internal/store/api/helper"
	categoryTypes "back/internal/store/api/types/category"
	materialTypes "back/internal/store/api/types/material"
	unitTypes "back/internal/store/api/types/unit"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"math"
	"net/http"
)

func Read(w http.ResponseWriter, r *http.Request, p materialTypes.T_params, url_q materialTypes.T_readQuery, q *pgstore.Queries) {
	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	offset := url_q.Page * url_q.PerPage
	limit := url_q.PerPage

	filter := ""

	type _response struct {
		Data []materialTypes.T_responseBody `json:"data"`
		Meta materialTypes.T_responseMeta   `json:"meta"`
	}

	is_first_field := true

	if url_q.FilterName != "" {
		if is_first_field {
			filter += " WHERE"
		} else {
			filter += " AND"
		}
		filter += " name ~* '" + url_q.FilterName + "'"
		is_first_field = false
	}

	if url_q.FilterUnitID != "" {
		if is_first_field {
			filter += " WHERE"
		} else {
			filter += " AND"
		}
		filter += " unit_id = '" + url_q.FilterUnitID + "'"
		is_first_field = false
	}

	if url_q.FilterCategoryID != "" {
		if is_first_field {
			filter += " WHERE"
		} else {
			filter += " AND"
		}
		filter += " category_id = '" + url_q.FilterCategoryID + "'"
		is_first_field = false
	}

	materials, err := q.C_FetchPaginatedMaterials(r.Context(), "SELECT id, name, description, quantity, category_id, unit_id FROM material"+filter+
		" ORDER BY "+url_q.SortColumn+" "+url_q.SortDirection+" LIMIT $1 OFFSET $2", pgstore.FetchPaginatedMaterialsParams{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		helper.HandleErrorMessage(w, err, "None material")
		return
	}

	size, err := q.C_GetTableSize(r.Context(), `SELECT count(*) AS exact_count FROM material`+filter)

	if err != nil {
		helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
		return
	}

	_materials := []materialTypes.T_responseBody{}
	for _, m := range materials {

		unit, err := q.FindUnitById(r.Context(), m.UnitID)

		if err != nil {
			helper.HandleErrorMessage(w, err, "Unit")
			return
		}

		category, err := q.FindCategoryById(r.Context(), m.CategoryID)

		if err != nil {
			helper.HandleErrorMessage(w, err, "Category")
			return
		}

		_materials = append(_materials, materialTypes.T_responseBody{
			ID:          m.ID.String(),
			Name:        m.Name,
			Description: m.Description.String,
			Quantity:    m.Quantity,
			Category: categoryTypes.T_responseBody{
				ID:   category.ID.String(),
				Name: category.Name,
			},
			Unit: unitTypes.T_responseBody{
				ID:        unit.ID.String(),
				Name:      unit.Name,
				ShortName: unit.ShortName.String,
			},
		})
	}

	total_pages := math.Ceil(float64(size) / float64(limit))

	data, _ := json.Marshal(_response{
		Data: _materials,
		Meta: materialTypes.T_responseMeta{
			Page:       url_q.Page,
			PerPage:    url_q.PerPage,
			TotalPages: int32(total_pages),
			Total:      int32(size),
		},
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
