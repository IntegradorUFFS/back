package categoryServices

import (
	"back/internal/store/api/helper"
	categoryTypes "back/internal/store/api/types/category"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func Create(w http.ResponseWriter, r *http.Request, p categoryTypes.T_params, b categoryTypes.T_body, q *pgstore.Queries) {

	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	category, err := q.CreateCategory(r.Context(), b.Name)

	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			helper.HandleError(w, "name", "already registered", http.StatusBadRequest)
		} else {
			helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
		}

		return
	}

	data, _ := json.Marshal(categoryTypes.T_responseWithMessage{
		Data: categoryTypes.T_responseBody{
			ID:   category.ID.String(),
			Name: category.Name,
		},
		Message: "Successfully created",
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
