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

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func Update(w http.ResponseWriter, r *http.Request, p materialTypes.T_params, b materialTypes.T_body, q *pgstore.Queries) {

	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	update_query := `UPDATE material set `

	query_end := ` WHERE id = $1 RETURNING name, description, quantity, category_id, unit_id`

	is_first_field := true

	if b.Name != "" {
		if !is_first_field {
			update_query = update_query + `, `
		}
		update_query = update_query + `name = '` + b.Name + `'`
		is_first_field = false
	}
	if b.Description != "" {
		if !is_first_field {
			update_query = update_query + `, `
		}
		update_query = update_query + `description = '` + b.Description + `'`
		is_first_field = false
	}
	if b.CategoryID != uuid.Nil {
		_, err := q.FindCategoryById(r.Context(), b.CategoryID)
		if err != nil {
			helper.HandleErrorMessage(w, err, "Category")
			return
		} else {
			if !is_first_field {
				update_query = update_query + `, `
			}
			update_query = update_query + `category_id = '` + b.CategoryID.String() + `'`
			is_first_field = false
		}
	}
	if b.UnitID != uuid.Nil {
		_, err := q.FindUnitById(r.Context(), b.UnitID)
		if err != nil {
			helper.HandleErrorMessage(w, err, "Unit")
			return
		} else {
			if !is_first_field {
				update_query = update_query + `, `
			}
			update_query = update_query + `unit_id = '` + b.UnitID.String() + `'`
			is_first_field = false
		}
	}

	material, err := q.C_UpdateMaterial(r.Context(), update_query+query_end, pgstore.UpdateMaterialParams{
		ID: p.ID,
	})

	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			helper.HandleError(w, "name", "already registered", http.StatusBadRequest)
		} else {
			helper.HandleErrorMessage(w, err, "Material")
		}

		return
	}

	c, err := q.FindCategoryById(r.Context(), material.CategoryID)
	if err != nil {
		helper.HandleErrorMessage(w, err, "Category")
		return
	}

	u, err := q.FindUnitById(r.Context(), material.UnitID)
	if err != nil {
		helper.HandleErrorMessage(w, err, "Unit")
		return
	}

	data, _ := json.Marshal(materialTypes.T_responseWithMessage{
		Data: materialTypes.T_responseBody{
			ID:          p.ID.String(),
			Name:        material.Name,
			Description: material.Description.String,
			Quantity:    material.Quantity,
			Category: categoryTypes.T_responseBody{
				ID:   c.ID.String(),
				Name: c.Name,
			},
			Unit: unitTypes.T_responseBody{
				ID:        u.ID.String(),
				Name:      u.Name,
				ShortName: u.ShortName.String,
			},
		},
		Message: "Successfully updated",
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
