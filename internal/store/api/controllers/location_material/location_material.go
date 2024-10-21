package locationMaterialController

import (
	"back/internal/store/api/helper"
	locationMaterialServices "back/internal/store/api/services/location_material"
	locationMaterialTypes "back/internal/store/api/types/location_material"
	pgstore "back/internal/store/pgstore/sqlc"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

type LocationMaterialQuery struct {
	q *pgstore.Queries
}

func (u LocationMaterialQuery) Find(w http.ResponseWriter, r *http.Request) {
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

	locationMaterialServices.Find(w, r, locationMaterialTypes.T_params{
		ID:            parsed_target_id,
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_requester_id,
	}, u.q)
}

func (u LocationMaterialQuery) List(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	url_query := locationMaterialTypes.T_readQuery{
		Page:          0,
		PerPage:       10,
		SortColumn:    "material.name",
		SortDirection: "ASC",
	}

	q_page := r.URL.Query().Get("page")
	q_per_page := r.URL.Query().Get("per_page")
	q_sort_column := r.URL.Query().Get("sort_column")
	q_sort_direction := r.URL.Query().Get("sort_direction")
	q_filter_location_id := r.URL.Query().Get("filter[location_id]")
	q_filter_material_id := r.URL.Query().Get("filter[material_id]")

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

	if q_sort_column == "id" || q_sort_column == "location.name" || q_sort_column == "material.name" {
		url_query.SortColumn = q_sort_column
	}

	if len(strings.TrimSpace(q_filter_location_id)) != 0 {
		parsed_unit_id, err := uuid.Parse(q_filter_location_id)
		if err == nil && parsed_unit_id != uuid.Nil {
			url_query.FilterLocationID = q_filter_location_id
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

	locationMaterialServices.Read(w, r, locationMaterialTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_id,
	}, url_query, u.q)
}

func (u LocationMaterialQuery) FindRelation(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	location_id := r.URL.Query().Get("location_id")
	material_id := r.URL.Query().Get("material_id")

	if claims["role"] == "viewer" || claims["id"] == "" {
		helper.HandleError(w, "", "Unauthorized user", http.StatusUnauthorized)
		return
	}

	parsed_id, err := uuid.Parse(claims["id"].(string))

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusUnprocessableEntity)
		return
	}

	if location_id == "" || len(strings.TrimSpace(location_id)) == 0 {
		helper.HandleError(w, "", "location_id query is missing", http.StatusUnprocessableEntity)
		return
	}

	if material_id == "" || len(strings.TrimSpace(material_id)) == 0 {
		helper.HandleError(w, "", "material_id query is missing", http.StatusUnprocessableEntity)
		return
	}

	parsed_location_id, err := uuid.Parse(location_id)

	if err != nil || parsed_location_id == uuid.Nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusUnprocessableEntity)
		return
	}

	parsed_material_id, err := uuid.Parse(material_id)

	if err != nil || parsed_material_id == uuid.Nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusUnprocessableEntity)
		return
	}

	locationMaterialServices.RelationFind(w, r, locationMaterialTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_id,
	}, locationMaterialTypes.T_relationQuery{
		LocationID: parsed_location_id,
		MaterialID: parsed_material_id,
	}, u.q)
}

func New(q *pgstore.Queries) LocationMaterialQuery {
	u := LocationMaterialQuery{
		q: q,
	}

	return u
}
