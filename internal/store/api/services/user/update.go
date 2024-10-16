package userServices

import (
	"back/internal/store/api/helper"
	userTypes "back/internal/store/api/types/user"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func Update(w http.ResponseWriter, r *http.Request, p userTypes.T_params, b userTypes.T_body, q *pgstore.Queries) {
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

	update_query := `UPDATE users
					  set `

	query_end := ` WHERE id = $1
				  RETURNING email, first_name, last_name, role`

	is_first_field := true

	if b.Email != "" {
		if !is_first_field {
			update_query = update_query + `, `
		}
		update_query = update_query + `email = '` + b.Email + `'`
		is_first_field = false
	}
	if b.FirstName != "" {
		if !is_first_field {
			update_query = update_query + `, `
		}
		update_query = update_query + `first_name = '` + b.FirstName + `'`
		is_first_field = false
	}
	if b.LastName != "" {
		if !is_first_field {
			update_query = update_query + `, `
		}
		update_query = update_query + `last_name = '` + b.LastName + `'`
		is_first_field = false
	}
	if b.Password != "" {
		if !is_first_field {
			update_query = update_query + `, `
		}
		update_query = update_query + `password = '` + b.Password + `'`
		is_first_field = false
	}
	if b.Role != "" {
		if !is_first_field {
			update_query = update_query + `, `
		}
		update_query = update_query + `role = '` + b.Role + `'`
		is_first_field = false
	}

	user_r, err := q.C_UpdateUser(r.Context(), update_query+query_end, pgstore.UpdateUserParams{
		ID: p.ID,
	})

	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			helper.HandleError(w, "email", "already registered", http.StatusBadRequest)
			return
		}
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	data, _ := json.Marshal(userTypes.T_responseWithMessage{
		Data: userTypes.T_responseBody{
			Email:     user_r.Email,
			FirstName: user_r.FirstName,
			LastName:  user_r.LastName.String,
			ID:        p.ID.String(),
			Role:      string(user_r.Role),
		},
		Message: "Successfully updated",
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)

}
