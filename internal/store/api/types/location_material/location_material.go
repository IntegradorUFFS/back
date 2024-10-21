package locationMaterialTypes

import (
	locationTypes "back/internal/store/api/types/location"
	materialTypes "back/internal/store/api/types/material"

	"github.com/google/uuid"
)

type T_jsonBody struct {
	Quantity   float32 `json:"quantity"`
	MaterialID string  `json:"material_id"`
	LocationID string  `json:"location_id"`
}

type T_body struct {
	Quantity   float32
	MaterialID uuid.UUID
	LocationID uuid.UUID
}

type T_params struct {
	ID            uuid.UUID
	RequesterID   uuid.UUID
	RequesterRole string
}

type T_responseBody struct {
	ID       string                       `json:"id"`
	Quantity float32                      `json:"quantity"`
	Material materialTypes.T_responseBody `json:"material"`
	Location locationTypes.T_responseBody `json:"location"`
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

type T_null struct {
}

type T_nullResponse struct {
	Data T_null `json:"data"`
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
	Page             int32
	PerPage          int32
	SortColumn       string
	SortDirection    string
	FilterMaterialID string
	FilterLocationID string
}

type T_relationQuery struct {
	LocationID uuid.UUID
	MaterialID uuid.UUID
}

type T_materialBody struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity float32 `json:"quantity"`
}

type T_updateBody struct {
	Quantity *float32 `json:"quantity"`
}
