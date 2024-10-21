package userServices

import (
	"back/internal/store/api/helper"
	userTypes "back/internal/store/api/types/user"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(w http.ResponseWriter, r *http.Request, b userTypes.T_body, q *pgstore.Queries) {

	user, err := q.FindUserByEmail(r.Context(), b.Email)

	if err != nil {
		helper.HandleError(w, "email", "Email not registered", http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(b.Password))

	if err != nil {
		helper.HandleError(w, "password", "Wrong password", http.StatusBadRequest)
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"id":   user.ID.String(),
		"role": user.Role,
	})

	if err != nil {
		helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(userTypes.T_responseWithMessageNToken{
		Data: userTypes.T_responseBodyNTOken{
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName.String,
			ID:        user.ID.String(),
			Role:      string(user.Role),
			Token:     tokenString,
		},
		Message: "Successfully authenticated",
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
