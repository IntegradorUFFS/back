package userTypes

import "github.com/google/uuid"

type T_responseBody struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ID        string `json:"id"`
	Role      string `json:"role"`
}

type T_responseBodyWithCreatedAt struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ID        string `json:"id"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type T_responseWithMessage struct {
	Data    T_responseBody `json:"data"`
	Message string         `json:"message"`
}

type T_responseBodyNTOken struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ID        string `json:"id"`
	Role      string `json:"role"`
	Token   string         `json:"token"`
}

type T_responseWithMessageNToken struct {
	Data    T_responseBodyNTOken `json:"data"`
	Message string         `json:"message"`
}

type T_response struct {
	Data T_responseBody `json:"data"`
}

type T_body struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ID        string `json:"id"`
	Role      string `json:"role"`
	Password  string `json:"password"`
}

type T_params struct {
	ID            uuid.UUID
	RequesterID   uuid.UUID
	RequesterRole string
}

type T_responseMessage struct {
	Message string `json:"message"`
}

type T_jsonMeta struct {
	Page    int32 `json:"page"`
	PerPage int32 `json:"per_page"`
}

type T_meta struct {
	Page    int32
	PerPage int32
}

type T_responseMeta struct {
	Page       int32 `json:"page"`
	PerPage    int32 `json:"per_page"`
	Total      int32 `json:"total"`
	TotalPages int32 `json:"total_pages"`
}

type T_readQuery struct {
	Page            int32
	PerPage         int32
	SortColumn      string
	SortDirection   string
	FilterFirstName string
	FilterLastName  string
	FilterEmail     string
	FilterRole      string
}
