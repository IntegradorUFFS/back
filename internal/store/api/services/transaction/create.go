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

func Create(w http.ResponseWriter, r *http.Request, p transactionTypes.T_params, b transactionTypes.T_body, q *pgstore.Queries) {
	requester_user, err := q.FindUserById(r.Context(), p.RequesterID)

	if err != nil || requester_user.Role != pgstore.Userrole(p.RequesterRole) {
		helper.HandleErrorMessage(w, err, "User")
		return
	}

	transaction_type := ""

	var origin pgstore.Location

	if b.OriginID != uuid.Nil {
		origin, err = q.FindLocationById(r.Context(), b.OriginID)

		transaction_type = "out"
		if err != nil {
			helper.HandleErrorMessage(w, err, "Origin")
			return
		}
	}

	var destiny pgstore.Location

	if b.DestinyID != uuid.Nil {
		destiny, err = q.FindLocationById(r.Context(), b.DestinyID)

		if transaction_type == "" {
			transaction_type = "in"
		} else {
			transaction_type = "transfer"
		}

		if err != nil {
			helper.HandleErrorMessage(w, err, "Destiny")
			return
		}
	}

	if transaction_type == "" {
		helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
		return
	}

	var origin_relation pgstore.LocationMaterial

	if b.OriginID != uuid.Nil {
		exists_origin, err := q.FindLocationMaterialByRelations(r.Context(), pgstore.FindLocationMaterialByRelationsParams{
			MaterialID: b.MaterialID,
			LocationID: b.OriginID,
		})

		if exists_origin.ID == uuid.Nil || err != nil {
			helper.HandleErrorMessage(w, err, "Origin")
			return
		}

		origin_relation = exists_origin
	}

	if (transaction_type == "transfer" || transaction_type == "out") && b.Quantity > origin_relation.Quantity {
		helper.HandleError(w, "quantity", "quantity higher than available quantity at origin", http.StatusBadRequest)
		return
	}

	var destiny_relation pgstore.LocationMaterial

	if b.DestinyID != uuid.Nil {
		exists_destiny, _ := q.FindLocationMaterialByRelations(r.Context(), pgstore.FindLocationMaterialByRelationsParams{
			MaterialID: b.MaterialID,
			LocationID: b.DestinyID,
		})

		if exists_destiny.ID != uuid.Nil {
			destiny_relation = exists_destiny
		}
	}

	if (transaction_type == "in" || transaction_type == "transfer") && destiny_relation.ID == uuid.Nil {
		destiny_relation, err = q.CreateLocationMaterial(r.Context(), pgstore.CreateLocationMaterialParams{
			Quantity:   b.Quantity,
			MaterialID: b.MaterialID,
			LocationID: b.DestinyID,
		})

		if err != nil {
			helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
			return
		}
	} else if (transaction_type == "in" || transaction_type == "transfer") && destiny_relation.ID != uuid.Nil {
		_, err = q.UpdateLocationMaterialQuantity(r.Context(), pgstore.UpdateLocationMaterialQuantityParams{
			ID:       destiny_relation.ID,
			Quantity: destiny_relation.Quantity + b.Quantity,
		})

		if err != nil {
			helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
			return
		}
	}

	material, err := q.FindMaterialById(r.Context(), b.MaterialID)

	if err != nil {
		helper.HandleErrorMessage(w, err, "Material")
		return
	}

	if transaction_type != "transfer" {
		quantity := b.Quantity
		if transaction_type == "out" {
			quantity *= -1
		}

		err = q.UpdateMaterialQuantity(r.Context(), pgstore.UpdateMaterialQuantityParams{
			ID:       b.MaterialID,
			Quantity: material.Quantity + quantity,
		})

		if err != nil {
			helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
			return
		}

		material.Quantity += quantity
	}

	if transaction_type == "out" || transaction_type == "transfer" {
		_, err = q.UpdateLocationMaterialQuantity(r.Context(), pgstore.UpdateLocationMaterialQuantityParams{
			ID:       origin_relation.ID,
			Quantity: origin_relation.Quantity - b.Quantity,
		})

		if err != nil {
			helper.HandleError(w, "", "Something went wrong", http.StatusInternalServerError)
			return
		}
	}

	var transaction pgstore.CreateTransactionRow

	if transaction_type == "transfer" {
		transaction, err = q.CreateTransaction(r.Context(), pgstore.CreateTransactionParams{
			Quantity:          b.Quantity,
			Type:              pgstore.Transactiontype(transaction_type),
			DestinyLocationID: destiny.ID,
			MaterialID:        material.ID,
			OriginLocationID:  origin.ID,
		})

		if err != nil {
			helper.HandleErrorMessage(w, err, "Transaction")
			return
		}
	} else if transaction_type == "in" {
		transaction, err = q.C_CreateTransactionWithDL(r.Context(), pgstore.CreateTransactionWithDLParams{
			Quantity:          b.Quantity,
			Type:              pgstore.Transactiontype(transaction_type),
			DestinyLocationID: destiny.ID,
			MaterialID:        material.ID,
		})

		if err != nil {
			helper.HandleErrorMessage(w, err, "Transaction")
			return
		}
	} else {
		transaction, err = q.C_CreateTransactionWithOL(r.Context(), pgstore.CreateTransactionWithOLParams{
			Quantity:         b.Quantity,
			Type:             pgstore.Transactiontype(transaction_type),
			OriginLocationID: origin.ID,
			MaterialID:       material.ID,
		})

		if err != nil {
			helper.HandleErrorMessage(w, err, "Transaction")
			return
		}
	}

	category, err := q.FindCategoryById(r.Context(), material.CategoryID)

	if err != nil {
		helper.HandleErrorMessage(w, err, "Category")
		return
	}

	unit, err := q.FindUnitById(r.Context(), material.UnitID)

	if err != nil {
		helper.HandleErrorMessage(w, err, "Unit")
		return
	}

	origin_data := &locationTypes.T_responseBody{}
	origin_data = nil

	if origin.ID != uuid.Nil {
		origin_data = &locationTypes.T_responseBody{
			ID:   origin.ID.String(),
			Name: origin.Name,
		}
	}

	destiny_data := &locationTypes.T_responseBody{}
	destiny_data = nil

	if destiny.ID != uuid.Nil {
		destiny_data = &locationTypes.T_responseBody{
			ID:   destiny.ID.String(),
			Name: destiny.Name,
		}
	}

	data, _ := json.Marshal(transactionTypes.T_responseWithMessage{
		Data: transactionTypes.T_responseBody{
			ID:       transaction.ID.String(),
			Quantity: transaction.Quantity,
			Material: materialTypes.T_responseBody{
				ID:          material.ID.String(),
				Name:        material.Name,
				Quantity:    material.Quantity,
				Description: material.Description.String,
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
			OriginLocation:  origin_data,
			DestinyLocation: destiny_data,
			CreatedAt:       transaction.CreatedAt.Time.Format(time.RFC3339),
			Type:            string(transaction.Type),
		},
		Message: "Transaction registered",
	})

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}
