package unitTypes

import (
	"net/url"

	"github.com/google/uuid"
)

type T_body struct {
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
}

type T_params struct {
	ID            uuid.UUID
	RequesterID   uuid.UUID
	RequesterRole string
}

type T_responseBody struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
}

type T_responseWithMessage struct {
	Data    T_responseBody `json:"data"`
	Message string         `json:"message"`
}

type T_responseMessage struct {
	Message string `json:"message"`
}

type T_response struct {
	Data T_responseBody `json:"data"`
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

type T_autocompleteQuery struct {
	ID     uuid.UUID
	Search string
}

type T_readQuery struct {
	Page    int32
	PerPage int32
	Query   url.Values
}
