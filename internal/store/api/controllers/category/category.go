package categoryController

import (
	"back/internal/store/api/helper"
	categoryServices "back/internal/store/api/services/category"
	categoryTypes "back/internal/store/api/types/category"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

type CategoryQuery struct {
	q *pgstore.Queries
}

func (u CategoryQuery) Create(w http.ResponseWriter, r *http.Request) {
	var body categoryTypes.T_body

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

	if body.Name == "" {
		helper.HandleError(w, "", "Some field is missing: name", http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(body.Name)) == 0 {
		helper.HandleError(w, "name", "Invalid input", http.StatusUnprocessableEntity)
		return
	}

	categoryServices.Create(w, r, categoryTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_id,
	}, body, u.q)
}

func (u CategoryQuery) Update(w http.ResponseWriter, r *http.Request) {
	var body categoryTypes.T_body
	id := chi.URLParam(r, "id")

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
		helper.HandleError(w, "", "Invalid uuid", http.StatusInternalServerError)
		return
	}

	categoryServices.Update(w, r, categoryTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_id,
		ID:            parsed_target_id,
	}, body, u.q)
}

func (u CategoryQuery) Delete(w http.ResponseWriter, r *http.Request) {
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
		helper.HandleError(w, "", "Invalid uuid", http.StatusInternalServerError)
		return
	}

	parsed_target_id, err := uuid.Parse(id)

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusInternalServerError)
		return
	}

	categoryServices.Delete(w, r, categoryTypes.T_params{
		ID:            parsed_target_id,
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_requester_id,
	}, u.q)
}

func (u CategoryQuery) Find(w http.ResponseWriter, r *http.Request) {
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
		helper.HandleError(w, "", "Invalid uuid", http.StatusInternalServerError)
		return
	}

	parsed_target_id, err := uuid.Parse(id)

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusInternalServerError)
		return
	}

	categoryServices.Find(w, r, categoryTypes.T_params{
		ID:            parsed_target_id,
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_requester_id,
	}, u.q)
}

func (u CategoryQuery) List(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	url_query := categoryTypes.T_readQuery{
		Page:    0,
		PerPage: 10,
	}

	q_page := r.URL.Query().Get("page")
	q_per_page := r.URL.Query().Get("per_page")

	if claims["role"] == "viewer" || claims["id"] == "" {
		helper.HandleError(w, "", "Unauthorized user", http.StatusUnauthorized)
		return
	}

	parsed_id, err := uuid.Parse(claims["id"].(string))

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusInternalServerError)
		return
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

	categoryServices.Read(w, r, categoryTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_id,
	}, url_query, u.q)
}

func (u CategoryQuery) Autocomplete(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	var url_query categoryTypes.T_autocompleteQuery

	q_search := r.URL.Query().Get("s")
	q_id := r.URL.Query().Get("id")

	if claims["role"] == "viewer" || claims["id"] == "" {
		helper.HandleError(w, "", "Unauthorized user", http.StatusUnauthorized)
		return
	}

	parsed_id, err := uuid.Parse(claims["id"].(string))

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusInternalServerError)
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

	categoryServices.Autocomplete(w, r, categoryTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_id,
	}, url_query, u.q)
}

func New(q *pgstore.Queries) CategoryQuery {
	u := CategoryQuery{
		q: q,
	}

	return u
}
