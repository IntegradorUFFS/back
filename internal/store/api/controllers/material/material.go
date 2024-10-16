package materialController

import (
	"back/internal/store/api/helper"
	materialServices "back/internal/store/api/services/material"
	materialTypes "back/internal/store/api/types/material"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

type MaterialQuery struct {
	q *pgstore.Queries
}

func (u MaterialQuery) Create(w http.ResponseWriter, r *http.Request) {
	var body materialTypes.T_jsonBody

	_, claims, _ := jwtauth.FromContext(r.Context())

	if claims["role"] == "viewer" || claims["id"] == "" {
		helper.HandleError(w, "", "Unauthorized user", http.StatusUnauthorized)
		return
	}

	parsed_id, err := uuid.Parse(claims["id"].(string))

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusInternalServerError)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		helper.HandleError(w, "", "Invalid json", http.StatusUnprocessableEntity)
		return
	}

	fields_err := []string{}

	if body.Name == "" {
		fields_err = append(fields_err, "name")
	}

	if body.CategoryID == "" {
		fields_err = append(fields_err, "category_id")
	}

	if body.UnitID == "" {
		fields_err = append(fields_err, "unit_id")
	}

	if len(fields_err) > 0 {
		helper.HandleError(w, "", "Some field is missing: "+strings.Join(fields_err, ", "), http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(body.Name)) == 0 {
		helper.HandleError(w, "name", "Invalid input", http.StatusUnprocessableEntity)
		return
	}

	if body.Description != "" && len(strings.TrimSpace(body.Description)) == 0 {
		helper.HandleError(w, "description", "Invalid input", http.StatusUnprocessableEntity)
		return
	}

	parsed_category_id, err := uuid.Parse(body.CategoryID)

	if err != nil || parsed_category_id == uuid.Nil {
		helper.HandleError(w, "category_id", "Invalid input", http.StatusUnprocessableEntity)
		return
	}

	parsed_unit_id, err := uuid.Parse(body.UnitID)

	if err != nil || parsed_unit_id == uuid.Nil {
		helper.HandleError(w, "unit_id", "Invalid input", http.StatusUnprocessableEntity)
		return
	}

	materialServices.Create(w, r, materialTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_id,
	}, materialTypes.T_body{
		Name:        body.Name,
		Description: body.Description,
		CategoryID:  parsed_category_id,
		UnitID:      parsed_unit_id,
	}, u.q)
}

func New(q *pgstore.Queries) MaterialQuery {
	u := MaterialQuery{
		q: q,
	}

	return u
}
