package userServices

import (
	"back/internal/store/api/helper"
	userTypes "back/internal/store/api/types/user"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/jwtauth/v5"
)

func Refresh(w http.ResponseWriter, r *http.Request, p userTypes.T_params, q *pgstore.Queries) {

	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"id":   requester_user.ID.String(),
		"role": requester_user.Role,
	})

	data, _ := json.Marshal(userTypes.T_responseWithMessageNToken{
		Data: userTypes.T_responseBodyNTOken{
			Email:     requester_user.Email,
			FirstName: requester_user.FirstName,
			LastName:  requester_user.LastName.String,
			ID:        requester_user.ID.String(),
			Role:      string(requester_user.Role),
			Token:     tokenString,
		},
		Message: "Successfully refreshed",
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
