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
	"net/http"

	"github.com/google/uuid"
)

func RelationFind(w http.ResponseWriter, r *http.Request, p locationMaterialTypes.T_params, url_q locationMaterialTypes.T_relationQuery, q *pgstore.Queries) {
	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	location_material, err := q.FindLocationMaterialByRelations(r.Context(), pgstore.FindLocationMaterialByRelationsParams{
		MaterialID: url_q.MaterialID,
		LocationID: url_q.LocationID,
	})

	if err != nil || location_material.ID == uuid.Nil {
		http.Error(w, "", http.StatusNotFound)
		data, _ := json.Marshal(locationMaterialTypes.T_nullResponse{})

		w.Header().Set("Content-Type", "application/json")

		w.Write(data)
		return
	}

	material, err := q.FindMaterialById(r.Context(), location_material.MaterialID)

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

	location, err := q.FindLocationById(r.Context(), location_material.LocationID)

	if err != nil {
		helper.HandleErrorMessage(w, err, "Location")
		return
	}

	data, _ := json.Marshal(locationMaterialTypes.T_response{
		Data: locationMaterialTypes.T_responseBody{
			ID:       location_material.ID.String(),
			Quantity: location_material.Quantity,
			Location: locationTypes.T_responseBody{
				ID:   location.ID.String(),
				Name: location.Name,
			},
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
		},
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
