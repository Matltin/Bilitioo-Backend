package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
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

type payPaymentResponse struct {
	Payment      db.Payment                `json:"payment"`
	Reservations []db.UpdateReservationRow `json:"reservations"`
	UserActivity db.UpdateUserActivityRow  `json:"user_activity"`
}

// payPayment godoc
//
//	@Summary		Pay for a reservation
//	@Description	Updates payment and reservation status. Requires authentication.
//	@Tags			Payment
//	@Accept			json
//	@Produce		json
//	@Param			request	body		payPaymentRequest	true	"Payment request body"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/payment [post]
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

	var amount int64 = 0

	log.Println(reservationsID)

	for _, r := range reservationsID {
		status, err := server.Queries.GetReservationStatus(ctx, r)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if status != db.TicketStatusRESERVING {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("reservation is not in reserving status, %s", status)))
			continue
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
			ID:     reservation.TicketID,
			Status: db.CheckReservationTicketStatus(req.ReservationStatus),
		}

		ticketAmount, err := server.Queries.UpdateTicketStatus(ctx, argTicket)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(errors.New("reserved not found")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// Invalidate ticket details cache
		ticketCacheKey := fmt.Sprintf("ticket_details:%d", reservation.ID)
		_ = server.redisClient.Delete(ctx, ticketCacheKey)

		// Invalidate all search ticket cache keys
		searchKeys, err := server.redisClient.Client.Keys(ctx, "search:*").Result()
		if err == nil {
			for _, key := range searchKeys {
				_ = server.redisClient.Delete(ctx, key)
			}
		}

		amount += ticketAmount

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

	if req.Type == "WALLET" {
		argWallet := db.AddToUserWalletParams{
			Wallet: -amount,
			UserID: authPayload.UserID,
		}

		err = server.Queries.AddToUserWallet(ctx, argWallet)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		user, err := server.Queries.GetUserByID(ctx, authPayload.UserID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		server.invalidateUserCache(ctx, authPayload.UserID, user.Email, user.PhoneNumber)

	}

	var userActivity db.UpdateUserActivityRow

	if req.UserActivityID > 0 {

		argUserActivity := db.UpdateUserActivityParams{
			ID:     req.UserActivityID,
			Status: db.ActivityStatusPURCHASED,
		}

		userActivity, err = server.Queries.UpdateUserActivity(ctx, argUserActivity)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, payPaymentResponse{
		Payment:      payment,
		Reservations: reservations,
		UserActivity: userActivity,
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
