package userServices

import (
	"back/internal/store/api/helper"
	userTypes "back/internal/store/api/types/user"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"net/http"
)

func Delete(w http.ResponseWriter, r *http.Request, p userTypes.T_params, q *pgstore.Queries) {

	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	user, err := q.FindUserById(r.Context(), p.ID)

	if err != nil {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	err = helper.HasUserCrudPermission(user, p)
	if err != nil {
		helper.HandleError(w, "", err.Error(), http.StatusUnauthorized)
		return
	}

	err = q.DeleteUser(r.Context(), p.ID)

	if err != nil {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	data, _ := json.Marshal(userTypes.T_responseMessage{
		Message: "Successfully deleted",
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)

}
