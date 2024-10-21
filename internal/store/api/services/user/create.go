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
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func Create(w http.ResponseWriter, r *http.Request, p userTypes.T_params, b userTypes.T_body, q *pgstore.Queries) {

	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(b.Password), bcrypt.DefaultCost)
	if err != nil {
		helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
		return
	}

	var user pgstore.CreateUserRow

	if b.Role != "" {
		user_r, err_r := q.CreateUserWithRole(r.Context(), pgstore.CreateUserWithRoleParams{
			Email:     b.Email,
			FirstName: b.FirstName,
			LastName:  pgtype.Text{String: b.LastName, Valid: true},
			Password:  string(hashed_password),
			Role:      pgstore.Userrole(b.Role),
		})

		user = pgstore.CreateUserRow(user_r)
		err = err_r

	} else {
		user, err = q.CreateUser(r.Context(), pgstore.CreateUserParams{
			Email:     b.Email,
			FirstName: b.FirstName,
			LastName:  pgtype.Text{String: b.LastName, Valid: true},
			Password:  string(hashed_password),
		})
	}

	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			helper.HandleError(w, "email", "already registered", http.StatusBadRequest)
		} else {
			helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
		}

		return
	}

	data, _ := json.Marshal(userTypes.T_responseWithMessage{
		Data: userTypes.T_responseBody{
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName.String,
			ID:        user.ID.String(),
			Role:      string(user.Role),
		},
		Message: "Successfully created",
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
