package unitServices

import (
	"back/internal/store/api/helper"
	unitTypes "back/internal/store/api/types/unit"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"net/http"
)

func Find(w http.ResponseWriter, r *http.Request, p unitTypes.T_params, q *pgstore.Queries) {

	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	unit, err := q.FindUnitById(r.Context(), p.ID)

	if err != nil {
		helper.HandleErrorMessage(w, err, "Unit")
		return
	}

	data, _ := json.Marshal(unitTypes.T_response{
		Data: unitTypes.T_responseBody{
			ShortName: unit.ShortName.String,
			Name:      unit.Name,
			ID:        unit.ID.String(),
		},
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
