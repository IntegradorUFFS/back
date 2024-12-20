package locationController

import (
	"back/internal/store/api/helper"
	locationServices "back/internal/store/api/services/location"
	locationTypes "back/internal/store/api/types/location"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

type LocationQuery struct {
	q *pgstore.Queries
}

func (u LocationQuery) Create(w http.ResponseWriter, r *http.Request) {
	var body locationTypes.T_body

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

	if body.Name == "" {
		helper.HandleError(w, "", "Some field is missing: name", http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(body.Name)) == 0 {
		helper.HandleError(w, "name", "Invalid input", http.StatusUnprocessableEntity)
		return
	}

	locationServices.Create(w, r, locationTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_id,
	}, body, u.q)
}

func (u LocationQuery) Update(w http.ResponseWriter, r *http.Request) {
	var body locationTypes.T_body
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

	if body.Name == "" {
		helper.HandleError(w, "", "Some field is missing: name", http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(body.Name)) == 0 {
		helper.HandleError(w, "name", "Invalid input", http.StatusUnprocessableEntity)
		return
	}

	parsed_target_id, err := uuid.Parse(id)

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusUnprocessableEntity)
		return
	}

	locationServices.Update(w, r, locationTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_id,
		ID:            parsed_target_id,
	}, body, u.q)
}

func (u LocationQuery) Delete(w http.ResponseWriter, r *http.Request) {
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

	locationServices.Delete(w, r, locationTypes.T_params{
		ID:            parsed_target_id,
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_requester_id,
	}, u.q)
}

func (u LocationQuery) Find(w http.ResponseWriter, r *http.Request) {
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

	locationServices.Find(w, r, locationTypes.T_params{
		ID:            parsed_target_id,
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_requester_id,
	}, u.q)
}

func (u LocationQuery) List(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	url_query := locationTypes.T_readQuery{
		Page:          0,
		PerPage:       10,
		SortColumn:    "name",
		SortDirection: "ASC",
		FilterName:    "",
	}

	q_page := r.URL.Query().Get("page")
	q_per_page := r.URL.Query().Get("per_page")
	q_sort_column := r.URL.Query().Get("sort_column")
	q_sort_direction := r.URL.Query().Get("sort_direction")
	q_filter_name := r.URL.Query().Get("filter[name]")

	if claims["role"] == "viewer" || claims["id"] == "" {
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

	if q_sort_column == "id" {
		url_query.SortColumn = q_sort_column
	}

	if len(strings.TrimSpace(q_filter_name)) != 0 {
		url_query.FilterName = q_filter_name
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

	locationServices.Read(w, r, locationTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_id,
	}, url_query, u.q)
}

func (u LocationQuery) Autocomplete(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	var url_query locationTypes.T_autocompleteQuery
	q_filter_material := r.URL.Query().Get("filter[material_id]")

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

	if len(strings.TrimSpace(q_filter_material)) != 0 {
		parsed_material_id, err := uuid.Parse(q_filter_material)
		if err == nil && parsed_material_id != uuid.Nil {
			url_query.FilterMaterialID = parsed_material_id
		}
	}

	parsed_query_id, err := uuid.Parse(q_id)

	if err == nil {
		url_query.ID = parsed_query_id
	}

	locationServices.Autocomplete(w, r, locationTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_id,
	}, url_query, u.q)
}

func New(q *pgstore.Queries) LocationQuery {
	u := LocationQuery{
		q: q,
	}

	return u
}
