package unitServices

import (
	"back/internal/store/api/helper"
	unitTypes "back/internal/store/api/types/unit"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"net/http"
)

func Delete(w http.ResponseWriter, r *http.Request, p unitTypes.T_params, q *pgstore.Queries) {

	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	err = q.DeleteUnit(r.Context(), p.ID)

	if err != nil {
		helper.HandleErrorMessage(w, err, "Unit")
		return
	}

	data, _ := json.Marshal(unitTypes.T_responseMessage{
		Message: "Successfully deleted",
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)

}
