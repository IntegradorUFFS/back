package locationServices

import (
	"back/internal/store/api/helper"
	locationTypes "back/internal/store/api/types/location"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func Update(w http.ResponseWriter, r *http.Request, p locationTypes.T_params, b locationTypes.T_body, q *pgstore.Queries) {

	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	location, err := q.UpdateLocation(r.Context(), pgstore.UpdateLocationParams{
		ID:   p.ID,
		Name: b.Name,
	})

	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			helper.HandleError(w, "name", "already registered", http.StatusBadRequest)
		} else {
			helper.HandleErrorMessage(w, err, "Location")
		}

		return
	}

	data, _ := json.Marshal(locationTypes.T_responseWithMessage{
		Data: locationTypes.T_responseBody{
			ID:   location.ID.String(),
			Name: location.Name,
		},
		Message: "Successfully updated",
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
