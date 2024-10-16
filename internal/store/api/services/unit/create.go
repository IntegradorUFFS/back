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
	"github.com/jackc/pgx/v5/pgtype"
)

func Create(w http.ResponseWriter, r *http.Request, p unitTypes.T_params, b unitTypes.T_body, q *pgstore.Queries) {

	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	unit, err := q.CreateUnit(r.Context(), pgstore.CreateUnitParams{
		Name:      b.Name,
		ShortName: pgtype.Text{String: b.ShortName, Valid: true},
	})

	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			helper.HandleError(w, "name", "already registered", http.StatusBadRequest)
		} else {
			helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
		}

		return
	}

	data, _ := json.Marshal(unitTypes.T_responseWithMessage{
		Data: unitTypes.T_responseBody{
			ID:        unit.ID.String(),
			Name:      unit.Name,
			ShortName: unit.ShortName.String,
		},
		Message: "Successfully created",
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
