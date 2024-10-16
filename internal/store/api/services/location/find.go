package locationServices

import (
	"back/internal/store/api/helper"
	locationTypes "back/internal/store/api/types/location"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"net/http"
)

func Find(w http.ResponseWriter, r *http.Request, p locationTypes.T_params, q *pgstore.Queries) {

	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	location, err := q.FindLocationById(r.Context(), p.ID)

	if err != nil {
		helper.HandleErrorMessage(w, err, "Location")
		return
	}

	data, _ := json.Marshal(locationTypes.T_response{
		Data: locationTypes.T_responseBody{
			Name: location.Name,
			ID:   location.ID.String(),
		},
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
