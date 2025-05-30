package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/token"
	"github.com/gin-gonic/gin"
)

type payPaymentRequest struct {
	PaymentID         int64  `json:"payment_id" binding:"required"`
	Type              string `json:"type" binding:"required,oneof=CASH CREDIT_CARD WALLET BANK_TRANSFER CRYPTO"`
	PaymentStatus     string `json:"payment_status" binding:"required,oneof=PENDING COMPLETED FAILED REFUNDED"`
	ReservationStatus string `json:"reservation_status" binding:"required,oneof=RESERVED RESERVING CANCELED CANCELED-BY-TIME"`
	UserActivityID    int64  `json:"user_activity_id"`
}

func (server *Server) payPayment(ctx *gin.Context) {
	var req payPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !isValidPaymentType(req.Type) {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid payment type")))
		return
	}

	if !isValidPaymentStatus(req.PaymentStatus) {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid payment status")))
		return
	}

	if !isValidTicketStatus(req.ReservationStatus) {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid reservation status")))
		return
	}

	argPayment := db.UpdatePaymentParams{
		Type:   db.PaymentType(req.Type),
		Status: db.PaymentStatus(req.PaymentStatus),
		ID:     req.PaymentID,
	}

	payment, err := server.Queries.UpdatePayment(ctx, argPayment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)
	var reservations []db.UpdateReservationRow

	reservationsID, err := server.Queries.GetIDReservation(context.Background(), req.PaymentID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	for _, r := range reservationsID {
		status, err := server.Queries.GetReservationStatus(ctx, r)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		arg := db.UpdateReservationParams{
			Status: db.TicketStatus(req.ReservationStatus),
			ID:     r,
		}

		reservation, err := server.Queries.UpdateReservation(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		reservations = append(reservations, reservation)

		argTicket := db.UpdateTicketStatusParams{
			ID:     reservation.ID,
			Status: db.CheckReservationTicketStatus(req.ReservationStatus),
		}

		err = server.Queries.UpdateTicketStatus(ctx, argTicket)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(errors.New("reserved not found")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		argChangeReservation := db.CreateChangeReservationParams{
			ReservationID: reservation.ID,
			AdminID:       1, // Assuming admin ID is 1, you can change this as needed
			UserID:        authPayload.UserID,
			FromStatus:    status,
			ToStatus:      db.TicketStatus(req.ReservationStatus),
		}

		_, err = server.Queries.CreateChangeReservation(ctx, argChangeReservation)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	var user_activity db.UpdateUserActivityRow

	if req.UserActivityID > 0 {

		argUserActivity := db.UpdateUserActivityParams{
			ID:     req.UserActivityID,
			Status: db.ActivityStatusPURCHASED,
		}

		user_activity, err = server.Queries.UpdateUserActivity(ctx, argUserActivity)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"payment":       payment,
		"reservations":  reservations,
		"user_activity": user_activity,
	})
}

func isValidPaymentType(t string) bool {
	switch t {
	case "CASH", "CREDIT_CARD", "WALLET", "BANK_TRANSFER", "CRYPTO":
		return true
	default:
		return false
	}
}

func isValidPaymentStatus(s string) bool {
	switch s {
	case "PENDING", "COMPLETED", "FAILED", "REFUNDED":
		return true
	default:
		return false
	}
}

func isValidTicketStatus(s string) bool {
	switch s {
	case "RESERVED", "RESERVING", "CANCELED", "CANCELED-BY-TIME":
		return true
	default:
		return false
	}
}
