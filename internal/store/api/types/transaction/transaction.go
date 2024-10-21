package transactionTypes

import (
	locationTypes "back/internal/store/api/types/location"
	materialTypes "back/internal/store/api/types/material"

	"github.com/google/uuid"
)

type T_params struct {
	ID            uuid.UUID
	RequesterID   uuid.UUID
	RequesterRole string
}

type T_locationResponseBody struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
}

type T_responseCleanBody struct {
	ID              string                            `json:"id"`
	Quantity        float32                           `json:"quantity"`
	CreatedAt       string                            `json:"created_at"`
	Type            string                            `json:"type"`
	Material        materialTypes.T_responseCleanBody `json:"material"`
	OriginLocation  *T_locationResponseBody           `json:"origin"`
	DestinyLocation *T_locationResponseBody           `json:"destiny"`
}

type T_responseBody struct {
	ID              string                        `json:"id"`
	Quantity        float32                       `json:"quantity"`
	CreatedAt       string                        `json:"created_at"`
	Type            string                        `json:"type"`
	Material        materialTypes.T_responseBody  `json:"material"`
	OriginLocation  *locationTypes.T_responseBody `json:"origin"`
	DestinyLocation *locationTypes.T_responseBody `json:"destiny"`
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
	Page                    int32
	PerPage                 int32
	SortColumn              string
	SortDirection           string
	FilterOriginLocationID  string
	FilterDestinyLocationID string
	FilterMaterialID        string
	FilterType              string
}

type T_jsonBody struct {
	Quantity   float32 `json:"quantity"`
	OriginID   string  `json:"origin_id"`
	DestinyID  string  `json:"destiny_id"`
	MaterialID string  `json:"material_id"`
}

type T_body struct {
	Quantity   float32
	OriginID   uuid.UUID
	DestinyID  uuid.UUID
	MaterialID uuid.UUID
}
