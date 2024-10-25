package unitServices

import (
	"back/internal/store/api/helper"
	unitTypes "back/internal/store/api/types/unit"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"errors"
	"fmt"
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
	if b.ShortName != "" {
		if update_count == 1 {
			update_query = update_query + `, `
		}
		update_query = update_query + `short_name = $` + fmt.Sprint(update_count)
		updates_arr = append(updates_arr, b.ShortName)
		update_count += 1
	}

	updates_arr = append(updates_arr, p.ID)

	unit, err := q.C_UpdateUnit(r.Context(), update_query+` WHERE id = $`+fmt.Sprint(update_count)+`
	RETURNING name, short_name`, updates_arr)

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
