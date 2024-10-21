package transactionController

import (
	"back/internal/store/api/helper"
	transactionServices "back/internal/store/api/services/transaction"
	transactionTypes "back/internal/store/api/types/transaction"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

type TransactionQuery struct {
	q *pgstore.Queries
}

func (u TransactionQuery) Find(w http.ResponseWriter, r *http.Request) {
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

	transactionServices.Find(w, r, transactionTypes.T_params{
		ID:            parsed_target_id,
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_requester_id,
	}, u.q)
}

func (u TransactionQuery) List(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	url_query := transactionTypes.T_readQuery{
		Page:          0,
		PerPage:       10,
		SortColumn:    "created_at",
		SortDirection: "DESC",
	}

	q_page := r.URL.Query().Get("page")
	q_per_page := r.URL.Query().Get("per_page")
	q_sort_column := r.URL.Query().Get("sort_column")
	q_sort_direction := r.URL.Query().Get("sort_direction")
	q_filter_type := r.URL.Query().Get("filter[type]")
	q_filter_origin_id := r.URL.Query().Get("filter[origin_id]")
	q_filter_destiny_id := r.URL.Query().Get("filter[destiny_id]")
	q_filter_material_id := r.URL.Query().Get("filter[material_id]")

	if claims["role"] == "viewer" || claims["id"] == "" {
		helper.HandleError(w, "", "Unauthorized user", http.StatusUnauthorized)
		return
	}

	parsed_id, err := uuid.Parse(claims["id"].(string))

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusUnprocessableEntity)
		return
	}

	if strings.ToLower(q_sort_direction) == "asc" {
		url_query.SortDirection = "ASC"
	}

	if q_sort_column == "id" || q_sort_column == "origin.name" || q_sort_column == "destiny.name" || q_sort_column == "material.name" || q_sort_column == "type" {
		url_query.SortColumn = q_sort_column
	}

	types_options := []string{"in", "out", "transfer"}

	if slices.Contains(types_options, q_filter_type) {
		url_query.FilterType = q_filter_type
	}

	if len(strings.TrimSpace(q_filter_origin_id)) != 0 {
		parsed_unit_id, err := uuid.Parse(q_filter_origin_id)
		if err == nil && parsed_unit_id != uuid.Nil {
			url_query.FilterOriginLocationID = q_filter_origin_id
		}
	}

	if len(strings.TrimSpace(q_filter_destiny_id)) != 0 {
		parsed_unit_id, err := uuid.Parse(q_filter_destiny_id)
		if err == nil && parsed_unit_id != uuid.Nil {
			url_query.FilterDestinyLocationID = q_filter_destiny_id
		}
	}

	if len(strings.TrimSpace(q_filter_material_id)) != 0 {
		parsed_category_id, err := uuid.Parse(q_filter_material_id)
		if err == nil && parsed_category_id != uuid.Nil {
			url_query.FilterMaterialID = q_filter_material_id
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

	transactionServices.Read(w, r, transactionTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_id,
	}, url_query, u.q)
}

func (u TransactionQuery) Create(w http.ResponseWriter, r *http.Request) {
	var body transactionTypes.T_jsonBody

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

	if body.DestinyID == "" && body.OriginID == "" {
		fields_err = append(fields_err, "destiny_id or origin_id")
	}

	if body.MaterialID == "" {
		fields_err = append(fields_err, "material_id")
	}

	if body.Quantity == 0 {
		fields_err = append(fields_err, "quantity (should be greater than 0)")
	}

	if len(fields_err) > 0 {
		helper.HandleError(w, "", "Some field is missing: "+strings.Join(fields_err, ", "), http.StatusBadRequest)
		return
	}

	if body.OriginID == body.DestinyID {
		helper.HandleError(w, "", "Origin and destiny can't be the same", http.StatusBadRequest)
		return
	}

	submit_body := transactionTypes.T_body{
		Quantity: body.Quantity,
	}

	parsed_origin_id, err := uuid.Parse(body.OriginID)

	if err == nil || parsed_origin_id != uuid.Nil {
		submit_body.OriginID = parsed_origin_id
	}

	parsed_destiny_id, err := uuid.Parse(body.DestinyID)

	if err == nil || parsed_destiny_id != uuid.Nil {
		submit_body.DestinyID = parsed_destiny_id
	}

	parsed_material_id, err := uuid.Parse(body.MaterialID)

	if err != nil || parsed_material_id == uuid.Nil {
		helper.HandleError(w, "material_id", "Invalid input", http.StatusUnprocessableEntity)
		return
	}

	submit_body.MaterialID = parsed_material_id

	transactionServices.Create(w, r, transactionTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_id,
	}, submit_body, u.q)
}

func New(q *pgstore.Queries) TransactionQuery {
	u := TransactionQuery{
		q: q,
	}

	return u
}
