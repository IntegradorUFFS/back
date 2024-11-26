package materialController

import (
	"back/internal/store/api/helper"
	materialServices "back/internal/store/api/services/material"
	materialTypes "back/internal/store/api/types/material"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
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
		helper.HandleError(w, "", "Invalid uuid", http.StatusUnprocessableEntity)
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

func (u MaterialQuery) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, claims, _ := jwtauth.FromContext(r.Context())

	if claims["role"] == "viewer" || claims["id"] == "" {
		helper.HandleError(w, "", "Unauthorized user", http.StatusUnauthorized)
		return
	}

	if id == "" || len(strings.TrimSpace(id)) == 0 {
		helper.HandleError(w, "", "Id param is missing", http.StatusUnprocessableEntity)
		return
	}

	parsed_requester_id, err := uuid.Parse(claims["id"].(string))

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusUnprocessableEntity)
		return
	}

	parsed_target_id, err := uuid.Parse(id)

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusUnprocessableEntity)
		return
	}

	materialServices.Delete(w, r, materialTypes.T_params{
		ID:            parsed_target_id,
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_requester_id,
	}, u.q)
}

func (u MaterialQuery) Find(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, claims, _ := jwtauth.FromContext(r.Context())

	if claims["role"] == "viewer" || claims["id"] == "" {
		helper.HandleError(w, "", "Unauthorized user", http.StatusUnauthorized)
		return
	}

	if id == "" || len(strings.TrimSpace(id)) == 0 {
		helper.HandleError(w, "", "Id param is missing", http.StatusUnprocessableEntity)
		return
	}

	parsed_requester_id, err := uuid.Parse(claims["id"].(string))

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusUnprocessableEntity)
		return
	}

	parsed_target_id, err := uuid.Parse(id)

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusUnprocessableEntity)
		return
	}

	materialServices.Find(w, r, materialTypes.T_params{
		ID:            parsed_target_id,
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_requester_id,
	}, u.q)
}

func (u MaterialQuery) Autocomplete(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	var url_query materialTypes.T_autocompleteQuery

	q_search := r.URL.Query().Get("s")
	q_id := r.URL.Query().Get("id")

	if claims["id"] == "" {
		helper.HandleError(w, "", "Unauthorized user", http.StatusUnauthorized)
		return
	}

	parsed_id, err := uuid.Parse(claims["id"].(string))

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusUnprocessableEntity)
		return
	}

	if len(strings.TrimSpace(q_search)) == 0 {
		url_query.Search = ""
	} else {
		url_query.Search = q_search
	}

	parsed_query_id, err := uuid.Parse(q_id)

	if err == nil {
		url_query.ID = parsed_query_id
	}

	materialServices.Autocomplete(w, r, materialTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_id,
	}, url_query, u.q)
}

func (u MaterialQuery) List(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	url_query := materialTypes.T_readQuery{
		Page:          0,
		PerPage:       10,
		SortColumn:    "name",
		SortDirection: "ASC",
	}

	q_page := r.URL.Query().Get("page")
	q_per_page := r.URL.Query().Get("per_page")
	q_sort_column := r.URL.Query().Get("sort_column")
	q_sort_direction := r.URL.Query().Get("sort_direction")
	q_filter_name := r.URL.Query().Get("filter[name]")
	q_filter_category_id := r.URL.Query().Get("filter[category_id]")
	q_filter_unit_id := r.URL.Query().Get("filter[unit_id]")

	if claims["id"] == "" {
		helper.HandleError(w, "", "Unauthorized user", http.StatusUnauthorized)
		return
	}

	parsed_id, err := uuid.Parse(claims["id"].(string))

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusUnprocessableEntity)
		return
	}

	if strings.ToLower(q_sort_direction) == "desc" {
		url_query.SortDirection = "DESC"
	}

	if q_sort_column == "id" || q_sort_column == "quantity" {
		url_query.SortColumn = q_sort_column
	}

	if len(strings.TrimSpace(q_filter_name)) != 0 {
		url_query.FilterName = q_filter_name
	}

	if len(strings.TrimSpace(q_filter_unit_id)) != 0 {
		parsed_unit_id, err := uuid.Parse(q_filter_unit_id)
		if err == nil && parsed_unit_id != uuid.Nil {
			url_query.FilterUnitID = q_filter_unit_id
		}
	}

	if len(strings.TrimSpace(q_filter_category_id)) != 0 {
		parsed_category_id, err := uuid.Parse(q_filter_category_id)
		if err == nil && parsed_category_id != uuid.Nil {
			url_query.FilterCategoryID = q_filter_category_id
		}
	}

	if q_page != "" && len(strings.TrimSpace(q_page)) > 0 {
		i, err := strconv.ParseInt(q_page, 10, 32)
		if err == nil {
			url_query.Page = int32(i)
		}
	}

	if q_per_page != "" && len(strings.TrimSpace(q_per_page)) > 0 {
		i, err := strconv.ParseInt(q_per_page, 10, 32)
		if err == nil && i > 0 {
			url_query.PerPage = int32(i)
		}
	}

	materialServices.Read(w, r, materialTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_id,
	}, url_query, u.q)
}

func (u MaterialQuery) Update(w http.ResponseWriter, r *http.Request) {
	var body materialTypes.T_jsonBody
	id := chi.URLParam(r, "id")

	_, claims, _ := jwtauth.FromContext(r.Context())

	if claims["role"] == "viewer" || claims["id"] == "" {
		helper.HandleError(w, "", "Unauthorized user", http.StatusUnauthorized)
		return
	}

	parsed_id, err := uuid.Parse(claims["id"].(string))

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusUnprocessableEntity)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		helper.HandleError(w, "", "Invalid json", http.StatusUnprocessableEntity)
		return
	}

	if body.Name == "" && body.Description == "" && body.CategoryID == "" && body.UnitID == "" {
		helper.HandleError(w, "", "At least one field is required: name, description, category_id, unit_id", http.StatusBadRequest)
		return
	}

	if body.Name != "" && len(strings.TrimSpace(body.Name)) == 0 {
		helper.HandleError(w, "name", "Invalid input", http.StatusUnprocessableEntity)
		return
	}

	if body.Description != "" && len(strings.TrimSpace(body.Description)) == 0 {
		helper.HandleError(w, "short_name", "Invalid input", http.StatusUnprocessableEntity)
		return
	}
	parsed_target_id, err := uuid.Parse(id)

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusUnprocessableEntity)
		return
	}

	submit_body := materialTypes.T_body{
		Name:        body.Name,
		Description: body.Description,
	}

	parsed_category_id, err := uuid.Parse(body.CategoryID)

	if body.CategoryID != "" && (err != nil || parsed_category_id == uuid.Nil) {
		helper.HandleError(w, "category_id", "Invalid input", http.StatusUnprocessableEntity)
		return
	}

	if err == nil {
		submit_body.CategoryID = parsed_category_id
	}

	parsed_unit_id, err := uuid.Parse(body.UnitID)

	if body.UnitID != "" && (err != nil || parsed_unit_id == uuid.Nil) {
		helper.HandleError(w, "unit_id", "Invalid input", http.StatusUnprocessableEntity)
		return
	}

	if err == nil {
		submit_body.UnitID = parsed_unit_id
	}

	materialServices.Update(w, r, materialTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_id,
		ID:            parsed_target_id,
	}, submit_body, u.q)
}

func New(q *pgstore.Queries) MaterialQuery {
	u := MaterialQuery{
		q: q,
	}

	return u
}
