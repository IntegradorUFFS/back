package transactionServices

import (
	"back/internal/store/api/helper"
	categoryTypes "back/internal/store/api/types/category"
	locationTypes "back/internal/store/api/types/location"
	materialTypes "back/internal/store/api/types/material"
	transactionTypes "back/internal/store/api/types/transaction"
	unitTypes "back/internal/store/api/types/unit"
	pgstore "back/internal/store/pgstore/sqlc"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func Find(w http.ResponseWriter, r *http.Request, p transactionTypes.T_params, q *pgstore.Queries) {

	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	transaction, err := q.FindTransactionById(r.Context(), p.ID)

	if err != nil {
		helper.HandleErrorMessage(w, err, "Transaction")
		return
	}

	destiny_data := &locationTypes.T_responseBody{}
	destiny_data = nil

	if transaction.DestinyLocationID != uuid.Nil {
		destiny, err := q.FindLocationById(r.Context(), transaction.DestinyLocationID)

		if err != nil {
			helper.HandleErrorMessage(w, err, "Location")
			return
		}

		destiny_data = &locationTypes.T_responseBody{
			ID:   destiny.ID.String(),
			Name: destiny.Name,
		}
	}

	origin_data := &locationTypes.T_responseBody{}
	origin_data = nil

	if transaction.OriginLocationID != uuid.Nil {
		origin, err := q.FindLocationById(r.Context(), transaction.OriginLocationID)

		if err != nil {
			helper.HandleErrorMessage(w, err, "Location")
			return
		}

		origin_data = &locationTypes.T_responseBody{
			ID:   origin.ID.String(),
			Name: origin.Name,
		}
	}

	material, err := q.FindMaterialById(r.Context(), transaction.MaterialID)

	if err != nil {
		helper.HandleErrorMessage(w, err, "Material")
		return
	}

	unit, err := q.FindUnitById(r.Context(), material.UnitID)

	if err != nil {
		helper.HandleErrorMessage(w, err, "Unit")
		return
	}

	category, err := q.FindCategoryById(r.Context(), material.CategoryID)

	if err != nil {
		helper.HandleErrorMessage(w, err, "Category")
		return
	}

	data, _ := json.Marshal(transactionTypes.T_response{
		Data: transactionTypes.T_responseBody{
			ID:       transaction.ID.String(),
			Quantity: transaction.Quantity,
			Material: materialTypes.T_responseBody{
				ID:          material.ID.String(),
				Name:        material.Name,
				Description: material.Description.String,
				Quantity:    material.Quantity,
				Category: categoryTypes.T_responseBody{
					ID:   category.ID.String(),
					Name: category.Name,
				},
				Unit: unitTypes.T_responseBody{
					ID:        unit.ID.String(),
					Name:      unit.Name,
					ShortName: unit.ShortName.String,
				},
			},
			CreatedAt:       transaction.CreatedAt.Time.Format(time.RFC3339),
			Type:            string(transaction.Type),
			OriginLocation:  origin_data,
			DestinyLocation: destiny_data,
		},
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
