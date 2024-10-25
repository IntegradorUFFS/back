package locationMaterialServices

import (
	"back/internal/store/api/helper"
	categoryTypes "back/internal/store/api/types/category"
	locationTypes "back/internal/store/api/types/location"
	locationMaterialTypes "back/internal/store/api/types/location_material"
	materialTypes "back/internal/store/api/types/material"
	unitTypes "back/internal/store/api/types/unit"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

func Read(w http.ResponseWriter, r *http.Request, p locationMaterialTypes.T_params, url_q locationMaterialTypes.T_readQuery, q *pgstore.Queries) {
	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	offset := url_q.Page * url_q.PerPage
	limit := url_q.PerPage

	filter := ""

	type _response struct {
		Data []locationMaterialTypes.T_responseBody `json:"data"`
		Meta locationMaterialTypes.T_responseMeta   `json:"meta"`
	}

	filter_count := 1

	var filters_arr []any

	if url_q.FilterLocationID != "" {
		if filter_count == 1 {
			filter += " WHERE"
		} else {
			filter += " AND"
		}
		filter += " location_id = $" + fmt.Sprint(filter_count)
		filter_count += 1

		filters_arr = append(filters_arr, url_q.FilterLocationID)
	}

	if url_q.FilterMaterialID != "" {
		if filter_count == 1 {
			filter += " WHERE"
		} else {
			filter += " AND"
		}
		filter += " material_id = $" + fmt.Sprint(filter_count)
		filter_count += 1

		filters_arr = append(filters_arr, url_q.FilterMaterialID)
	}

	size_filters := filters_arr

	filters_arr = append(filters_arr, limit)
	filters_arr = append(filters_arr, offset)

	location_materials, err := q.C_FetchPaginatedLocationMaterials(r.Context(), `SELECT location_material.*
		FROM location_material
		LEFT JOIN material ON location_material.material_id=material.id
		LEFT JOIN location ON location_material.location_id=location.id`+filter+
		" ORDER BY "+url_q.SortColumn+" "+url_q.SortDirection+" LIMIT $" + fmt.Sprint(filter_count) + " OFFSET $" + fmt.Sprint(filter_count + 1), filters_arr)

	if err != nil {
		helper.HandleErrorMessage(w, err, "None material")
		return
	}

	size, err := q.C_GetTableSize(r.Context(), `SELECT count(*) AS exact_count FROM location_material`+filter, size_filters)

	if err != nil {
		helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
		return
	}

	_location_materials := []locationMaterialTypes.T_responseBody{}
	for _, lm := range location_materials {

		material, err := q.FindMaterialById(r.Context(), lm.MaterialID)

		if err != nil {
			helper.HandleErrorMessage(w, err, "Material")
			return
		}

		category, err := q.FindCategoryById(r.Context(), material.CategoryID)

		if err != nil {
			helper.HandleErrorMessage(w, err, "Category")
			return
		}

		unit, err := q.FindUnitById(r.Context(), material.UnitID)

		if err != nil {
			helper.HandleErrorMessage(w, err, "Unit")
			return
		}

		location, err := q.FindLocationById(r.Context(), lm.LocationID)

		if err != nil {
			helper.HandleErrorMessage(w, err, "Location")
			return
		}

		_location_materials = append(_location_materials, locationMaterialTypes.T_responseBody{
			ID:       lm.ID.String(),
			Quantity: lm.Quantity,
			Material: materialTypes.T_responseBody{
				ID:          material.ID.String(),
				Name:        material.Name,
				Description: material.Description.String,
				Quantity:    material.Quantity,
				Category: categoryTypes.T_responseBody{
					ID:   category.ID.String(),
					Name: category.Name,
				},
				Unit: unitTypes.T_responseBody{
					ID:        unit.ID.String(),
					Name:      unit.Name,
					ShortName: unit.ShortName.String,
				},
			},
			Location: locationTypes.T_responseBody{
				ID:   location.ID.String(),
				Name: location.Name,
			},
		})
	}

	total_pages := math.Ceil(float64(size) / float64(limit))

	data, _ := json.Marshal(_response{
		Data: _location_materials,
		Meta: locationMaterialTypes.T_responseMeta{
			Page:       url_q.Page,
			PerPage:    url_q.PerPage,
			TotalPages: int32(total_pages),
			Total:      int32(size),
		},
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
