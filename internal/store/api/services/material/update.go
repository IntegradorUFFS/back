package materialServices

import (
	"back/internal/store/api/helper"
	categoryTypes "back/internal/store/api/types/category"
	materialTypes "back/internal/store/api/types/material"
	unitTypes "back/internal/store/api/types/unit"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"errors"
	"fmt"
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

	update_count := 1

	var updates_arr []any

	if b.Name != "" {
		if update_count == 1 {
			update_query = update_query + `, `
		}
		update_query = update_query + `name = $` + fmt.Sprint(update_count)
		updates_arr = append(updates_arr, b.Name)
		update_count += 1
	}
	if b.Description != "" {
		if update_count == 1 {
			update_query = update_query + `, `
		}
		update_query = update_query + `description = $` + fmt.Sprint(update_count)
		updates_arr = append(updates_arr, b.Description)
		update_count += 1
	}
	if b.CategoryID != uuid.Nil {
		_, err := q.FindCategoryById(r.Context(), b.CategoryID)
		if err != nil {
			helper.HandleErrorMessage(w, err, "Category")
			return
		} else {
			if update_count == 1 {
				update_query = update_query + `, `
			}
			update_query = update_query + `category_id = $` + fmt.Sprint(update_count)
			updates_arr = append(updates_arr, b.CategoryID.String())
			update_count += 1
		}
	}
	if b.UnitID != uuid.Nil {
		_, err := q.FindUnitById(r.Context(), b.UnitID)
		if err != nil {
			helper.HandleErrorMessage(w, err, "Unit")
			return
		} else {
			if update_count == 1 {
				update_query = update_query + `, `
			}
			update_query = update_query + `unit_id = $` + fmt.Sprint(update_count)
			updates_arr = append(updates_arr, b.UnitID.String())
			update_count += 1
		}
	}

	updates_arr = append(updates_arr, p.ID)

	material, err := q.C_UpdateMaterial(r.Context(), update_query+` WHERE id = $`+fmt.Sprint(update_count)+` RETURNING name, description, quantity, category_id, unit_id`, updates_arr)

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
