package unitServices

import (
	"back/internal/store/api/helper"
	unitTypes "back/internal/store/api/types/unit"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func Update(w http.ResponseWriter, r *http.Request, p unitTypes.T_params, b unitTypes.T_body, q *pgstore.Queries) {

	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	update_query := `UPDATE unit
					  set `

	query_end := ` WHERE id = $1
				  RETURNING name, short_name`

	is_first_field := true

	if b.Name != "" {
		if !is_first_field {
			update_query = update_query + `, `
		}
		update_query = update_query + `name = '` + b.Name + `'`
		is_first_field = false
	}
	if b.ShortName != "" {
		if !is_first_field {
			update_query = update_query + `, `
		}
		update_query = update_query + `short_name = '` + b.ShortName + `'`
		is_first_field = false
	}

	unit, err := q.C_UpdateUnit(r.Context(), update_query+query_end, pgstore.UpdateUnitParams{
		ID: p.ID,
	})

	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			helper.HandleError(w, "name", "already registered", http.StatusBadRequest)
		} else {
			helper.HandleErrorMessage(w, err, "Unit")
		}

		return
	}

	data, _ := json.Marshal(unitTypes.T_responseWithMessage{
		Data: unitTypes.T_responseBody{
			ID:        p.ID.String(),
			Name:      unit.Name,
			ShortName: unit.ShortName.String,
		},
		Message: "Successfully updated",
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
