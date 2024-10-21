package userController

import (
	"back/internal/store/api/helper"
	userServices "back/internal/store/api/services/user"
	userTypes "back/internal/store/api/types/user"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserQuery struct {
	q *pgstore.Queries
}

func (u UserQuery) Create(w http.ResponseWriter, r *http.Request) {
	var body userTypes.T_body

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

	if body.Email == "" {
		fields_err = append(fields_err, "email")
	}
	if body.FirstName == "" {
		fields_err = append(fields_err, "first_name")
	}

	if body.Password == "" {
		fields_err = append(fields_err, "password")
	}

	if len(fields_err) > 0 {
		helper.HandleError(w, "", "Some field is missing: "+strings.Join(fields_err, ", "), http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(body.Email)) == 0 {
		helper.HandleError(w, "email", "Invalid input", http.StatusUnprocessableEntity)
		return
	}
	if len(strings.TrimSpace(body.FirstName)) == 0 {
		helper.HandleError(w, "first_name", "Invalid input", http.StatusUnprocessableEntity)
		return
	}
	if body.LastName != "" && len(strings.TrimSpace(body.LastName)) == 0 {
		helper.HandleError(w, "last_name", "Invalid input", http.StatusUnprocessableEntity)
		return
	}

	if len(strings.TrimSpace(body.Password)) == 0 {
		helper.HandleError(w, "password", "Invalid input", http.StatusUnprocessableEntity)
		return
	}

	err = checkmail.ValidateFormat(body.Email)
	if err != nil {
		helper.HandleError(w, "email", "Invalid email", http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(body.Role)) == 0 || body.Role == "" || claims["role"] == "manager" {
		userServices.Create(w, r, userTypes.T_params{
			RequesterRole: claims["role"].(string),
			RequesterID:   parsed_id,
		}, userTypes.T_body{
			Email:     body.Email,
			FirstName: body.FirstName,
			LastName:  body.LastName,
			Password:  body.Password,
			Role:      "",
		}, u.q)
		return
	}

	if body.Role != "viewer" && body.Role != "manager" {
		helper.HandleError(w, "role", "Invalid role", http.StatusBadRequest)
		return
	}

	userServices.Create(w, r, userTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_id,
	}, body, u.q)
}

func (u UserQuery) Auth(w http.ResponseWriter, r *http.Request) {
	var body userTypes.T_body

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		helper.HandleError(w, "", "Invalid json", http.StatusUnprocessableEntity)
		return
	}

	fields := []string{}

	if body.Email == "" {
		fields = append(fields, "email")
	}

	if body.Password == "" {
		fields = append(fields, "password")
	}

	if len(fields) > 0 {
		helper.HandleError(w, "", "Some field is missing: "+strings.Join(fields, ", "), http.StatusUnprocessableEntity)
		return
	}

	if len(strings.TrimSpace(body.Email)) == 0 {
		helper.HandleError(w, "email", "Invalid input", http.StatusUnprocessableEntity)
		return
	}

	if len(strings.TrimSpace(body.Password)) == 0 {
		helper.HandleError(w, "password", "Invalid input", http.StatusUnprocessableEntity)
		return
	}

	err := checkmail.ValidateFormat(body.Email)
	if err != nil {
		helper.HandleError(w, "email", "Invalid email", http.StatusBadRequest)
		return
	}

	userServices.Authenticate(w, r, body, u.q)
}

func (u UserQuery) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, claims, _ := jwtauth.FromContext(r.Context())

	if (claims["role"] == "viewer" && id != claims["id"]) || claims["id"] == "" {
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

	userServices.Delete(w, r, userTypes.T_params{
		ID:            parsed_target_id,
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_requester_id,
	}, u.q)
}

func (u UserQuery) Find(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, claims, _ := jwtauth.FromContext(r.Context())

	if (claims["role"] == "viewer" && id != claims["id"]) || claims["id"] == "" {
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

	parsed_id, err := uuid.Parse(id)

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusUnprocessableEntity)
		return
	}

	userServices.Find(w, r, userTypes.T_params{
		ID:            parsed_id,
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_requester_id,
	}, u.q)
}

func (u UserQuery) List(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	url_query := userTypes.T_readQuery{
		Page:          0,
		PerPage:       10,
		SortColumn:    "first_name",
		SortDirection: "ASC",
	}

	q_page := r.URL.Query().Get("page")
	q_per_page := r.URL.Query().Get("per_page")
	q_sort_column := r.URL.Query().Get("sort_column")
	q_sort_direction := r.URL.Query().Get("sort_direction")
	q_filter_first_name := r.URL.Query().Get("filter[first_name]")
	q_filter_last_name := r.URL.Query().Get("filter[last_name]")
	q_filter_role := r.URL.Query().Get("filter[role]")
	q_filter_email := r.URL.Query().Get("filter[email]")

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

	possible_fields := []string{"id", "last_name", "email", "role"}

	if slices.Contains(possible_fields, q_sort_column) {
		url_query.SortColumn = q_sort_column
	}

	if len(strings.TrimSpace(q_filter_first_name)) != 0 {
		url_query.FilterFirstName = q_filter_first_name
	}

	if len(strings.TrimSpace(q_filter_last_name)) != 0 {
		url_query.FilterLastName = q_filter_last_name
	}

	if len(strings.TrimSpace(q_filter_email)) != 0 {
		url_query.FilterEmail = q_filter_email
	}

	if len(strings.TrimSpace(q_filter_role)) != 0 {
		url_query.FilterRole = q_filter_role
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

	userServices.Read(w, r, userTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_id,
	}, url_query, u.q)
}

func (u UserQuery) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body userTypes.T_body

	_, claims, _ := jwtauth.FromContext(r.Context())

	if (claims["role"] == "viewer" && id != claims["id"]) || claims["id"] == "" {
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

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		helper.HandleError(w, "", "Invalid json", http.StatusUnprocessableEntity)
		return
	}

	if body.Email == "" && body.FirstName == "" && body.LastName == "" && body.Password == "" && body.Role == "" {
		helper.HandleError(w, "", "At least one field is required: email, first_name, last_name, password, role", http.StatusBadRequest)
		return
	}

	if body.Email != "" && len(strings.TrimSpace(body.Email)) == 0 {
		helper.HandleError(w, "email", "Invalid input", http.StatusUnprocessableEntity)
		return
	}
	if body.FirstName != "" && len(strings.TrimSpace(body.FirstName)) == 0 {
		helper.HandleError(w, "first_name", "Invalid input", http.StatusUnprocessableEntity)
		return
	}
	if body.Password != "" && len(strings.TrimSpace(body.Password)) == 0 {
		helper.HandleError(w, "password", "Invalid input", http.StatusUnprocessableEntity)
		return
	}
	if body.LastName != "" && len(strings.TrimSpace(body.LastName)) == 0 {
		helper.HandleError(w, "last_name", "Invalid input", http.StatusUnprocessableEntity)
		return
	}

	parsed_id, err := uuid.Parse(id)

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusUnprocessableEntity)
		return
	}

	var target_user = userTypes.T_body{
		Email:     "",
		FirstName: "",
		LastName:  "",
		Password:  "",
	}

	if body.Email != "" {
		err = checkmail.ValidateFormat(body.Email)
		if err != nil {
			helper.HandleError(w, "email", "Invalid email", http.StatusBadRequest)
			return
		}
		target_user.Email = body.Email
	}
	if body.FirstName != "" {
		target_user.FirstName = body.FirstName
	}
	if body.LastName != "" {
		target_user.LastName = body.LastName
	}
	if body.Password != "" {
		hashed_password, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
			return
		}

		target_user.Password = string(hashed_password)
	}

	if len(strings.TrimSpace(body.Role)) == 0 || body.Role == "" || claims["role"].(string) == "manager" {
		userServices.Update(w, r, userTypes.T_params{
			ID:            parsed_id,
			RequesterRole: claims["role"].(string),
			RequesterID:   parsed_requester_id,
		}, userTypes.T_body{
			Email:     target_user.Email,
			FirstName: target_user.FirstName,
			LastName:  target_user.LastName,
			Password:  target_user.Password,
			Role:      "",
		}, u.q)
		return
	}

	if body.Role != "viewer" && body.Role != "manager" {
		helper.HandleError(w, "role", "Invalid role", http.StatusBadRequest)
		return
	}

	target_user.Role = body.Role

	userServices.Update(w, r, userTypes.T_params{
		ID:            parsed_id,
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_requester_id,
	}, target_user, u.q)
}

func (u UserQuery) Refresh(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	parsed_requester_id, err := uuid.Parse(claims["id"].(string))

	if err != nil {
		helper.HandleError(w, "", "Invalid uuid", http.StatusUnprocessableEntity)
		return
	}

	userServices.Refresh(w, r, userTypes.T_params{
		RequesterRole: claims["role"].(string),
		RequesterID:   parsed_requester_id,
	}, u.q)
}

func New(q *pgstore.Queries) UserQuery {
	u := UserQuery{
		q: q,
	}

	return u
}
