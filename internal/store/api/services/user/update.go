package userServices

import (
	"back/internal/store/api/helper"
	userTypes "back/internal/store/api/types/user"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"errors"
	"fmt"
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

	update_count := 1

	var updates_arr []any

	if b.FirstName != "" {
		if update_count > 1 {
			update_query = update_query + `, `
		}
		update_query = update_query + `first_name = $` + fmt.Sprint(update_count)
		updates_arr = append(updates_arr, b.FirstName)
		update_count += 1
	}

	if b.Email != "" {
		if update_count > 1 {
			update_query = update_query + `, `
		}
		update_query = update_query + `email = $` + fmt.Sprint(update_count)
		updates_arr = append(updates_arr, b.Email)
		update_count += 1
	}

	if b.LastName != "" {
		if update_count > 1 {
			update_query = update_query + `, `
		}
		update_query = update_query + `last_name = $` + fmt.Sprint(update_count)
		updates_arr = append(updates_arr, b.LastName)
		update_count += 1
	}

	if b.Password != "" {
		if update_count > 1 {
			update_query = update_query + `, `
		}
		update_query = update_query + `password = $` + fmt.Sprint(update_count)
		updates_arr = append(updates_arr, b.Password)
		update_count += 1
	}

	if b.Role != "" {
		if update_count > 1 {
			update_query = update_query + `, `
		}
		update_query = update_query + `role = $` + fmt.Sprint(update_count)
		updates_arr = append(updates_arr, b.Role)
		update_count += 1
	}

	updates_arr = append(updates_arr, p.ID)

	user_r, err := q.C_UpdateUser(r.Context(), update_query+` WHERE id = $`+fmt.Sprint(update_count)+`
	RETURNING email, first_name, last_name, role`, updates_arr)

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
