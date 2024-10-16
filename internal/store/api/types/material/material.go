package materialTypes

import (
	categoryTypes "back/internal/store/api/types/category"
	unitTypes "back/internal/store/api/types/unit"

	"github.com/google/uuid"
)

type T_jsonBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryID  string `json:"category_id"`
	UnitID      string `json:"unit_id"`
}

type T_body struct {
	Name        string
	Description string
	CategoryID  uuid.UUID
	UnitID      uuid.UUID
}

type T_params struct {
	ID            uuid.UUID
	RequesterID   uuid.UUID
	RequesterRole string
}

type T_responseBody struct {
	ID          string                       `json:"id"`
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Quantity    float32                      `json:"quantity"`
	Category    categoryTypes.T_responseBody `json:"category"`
	Unit        unitTypes.T_responseBody     `json:"unit"`
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
	Page       int64 `json:"page"`
	PerPage    int64 `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

type T_autocompleteQuery struct {
	ID     uuid.UUID
	Search string
}

type T_readQuery struct {
	Page    int32
	PerPage int32
}
