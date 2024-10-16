package helper

import (
	userTypes "back/internal/store/api/types/user"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type _fieldErrResponse struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type _errResponse struct {
	Error string `json:"error"`
}

func HandleError(w http.ResponseWriter, f string, e string, code int) {
	if f != "" {
		http.Error(w, "", code)
		data, _ := json.Marshal(_fieldErrResponse{
			Field: f,
			Error: e,
		})

		w.Header().Set("Content-Type", "application/json")

		w.Write(data)
	} else {
		http.Error(w, "", code)
		data, _ := json.Marshal(_errResponse{
			Error: e,
		})

		w.Header().Set("Content-Type", "application/json")

		w.Write(data)
	}
}

func HasUserCrudPermission(u pgstore.FindUserByIdRow, p userTypes.T_params) error {

	if p.RequesterRole == "viewer" {
		if p.RequesterID != u.ID {
			return errors.New("Insufficient permissions")
		}
	}

	if p.RequesterRole == "manager" {
		if u.Role == "admin" {
			return errors.New("Insufficient permissions")
		}
		if u.Role == "manager" && p.RequesterID != u.ID {
			return errors.New("Insufficient permissions")
		}
	}

	return nil
}

func HandleErrorMessage(w http.ResponseWriter, err error, origin string) {
	if errors.Is(err, pgx.ErrNoRows) {
		HandleError(w, "", origin+" not registered", http.StatusNotFound)
		return
	}

	var e *pgconn.PgError
	if errors.As(err, &e) && e.Code == pgerrcode.ForeignKeyViolation {
		HandleError(w, "", origin+" in use by some other table", http.StatusBadRequest)
		return
	}

	// if errors.Is(err, pgx.ERR) {
	// 	HandleError(w, "", origin+" not registered", http.StatusNotFound)
	// 	return
	// }

	HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
}
