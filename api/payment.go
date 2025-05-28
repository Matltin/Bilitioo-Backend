package api

import (
	"fmt"
	"net/http"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/token"
	"github.com/gin-gonic/gin"
)

type payPaymentRequest struct {
	PaymentID         int64   `json:"payment_id"`
	Reservations      []int64 `json:"reservations"`
	Type              string  `json:"type"`
	PaymentStatus     string  `json:"payment_status"`
	ReservationStatus string  `json:"reservatoin_status"`
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
	var reservations []db.Reservation

	for _, r := range req.Reservations {
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

	userActivityID := ctx.MustGet(userActivityID).(int64)

	argUserActivity := db.UpdateUserActivityParams{
		ID:     userActivityID,
		Status: db.ActivityStatusPURCHASED,
	}

	server.Queries.UpdateUserActivity(ctx, argUserActivity)

	ctx.JSON(http.StatusOK, gin.H{
		"payment":          payment,
		"reservations":     reservations,
		"user_activity_id": userActivityID,
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
