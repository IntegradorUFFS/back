package materialServices

import (
	"back/internal/store/api/helper"
	categoryTypes "back/internal/store/api/types/category"
	materialTypes "back/internal/store/api/types/material"
	unitTypes "back/internal/store/api/types/unit"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func Create(w http.ResponseWriter, r *http.Request, p materialTypes.T_params, b materialTypes.T_body, q *pgstore.Queries) {

	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	category, err := q.FindCategoryById(r.Context(), b.CategoryID)

	if err != nil {
		helper.HandleErrorMessage(w, err, "Category")
		return
	}

	unit, err := q.FindUnitById(r.Context(), b.UnitID)

	if err != nil {
		helper.HandleErrorMessage(w, err, "Unit")
		return
	}

	material, err := q.CreateMaterial(r.Context(), pgstore.CreateMaterialParams{
		Name:        b.Name,
		Description: pgtype.Text{String: b.Description, Valid: true},
		CategoryID:  b.CategoryID,
		UnitID:      b.UnitID,
	})

	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			helper.HandleError(w, "name", "already registered", http.StatusBadRequest)
		} else {
			helper.HandleError(w, "", err.Error(), http.StatusInternalServerError)
		}

		return
	}

	data, _ := json.Marshal(materialTypes.T_responseWithMessage{
		Data: materialTypes.T_responseBody{
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
		Message: "Successfully created",
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
