package materialServices

import (
	"back/internal/store/api/helper"
	categoryTypes "back/internal/store/api/types/category"
	materialTypes "back/internal/store/api/types/material"
	unitTypes "back/internal/store/api/types/unit"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"net/http"
)

func Find(w http.ResponseWriter, r *http.Request, p materialTypes.T_params, q *pgstore.Queries) {

	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	material, err := q.FindMaterialById(r.Context(), p.ID)

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

	data, _ := json.Marshal(materialTypes.T_response{
		Data: materialTypes.T_responseBody{
			Name:        material.Name,
			ID:          material.ID.String(),
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
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
