package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type getTicketPenaltiesRequest struct {
	TicketID int64 `uri:"ticket_id" binding:"required"`
}

// getTicketPenalties godoc
//
//	@Summary		Get penalties for a ticket
//	@Description	Retrieve penalties associated with a specific ticket. Requires authentication.
//	@Tags			Penalties
//	@Accept			json
//	@Produce		json
//	@Param			ticket_id	path		int			true	"Ticket ID"
//	@Success		200			{array}		db.Penalty	"List of penalties"
//	@Failure		400			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/ticket-penalties/{ticket_id} [get]
func (server *Server) getTicketPenalties(ctx *gin.Context) {
	var req getTicketPenaltiesRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	penalties, err := server.Queries.GetTicketPenalties(ctx, req.TicketID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, penalties)

}

type cancelReservationRequest struct {
	TicketID int64 `uri:"ticket_id" binding:"required,min=1"`
}

type cancelReservationResponse struct {
	Message           string               `json:"message"`
	AmountRefunded    int64                `json:"amount_refunded"`
	TicketID          int64                `json:"ticket_id"`
	ChangeReservation db.ChangeReservation `json:"change_reservation"`
}

// cancelReservation godoc
//
//	@Summary		Cancel a reservation
//	@Description	Cancel a ticket reservation and calculate penalty refund. Requires authentication.
//	@Tags			Reservation
//	@Accept			json
//	@Produce		json
//	@Param			ticket_id	path		int	true	"Ticket ID"
//	@Success		200			{object}	cancelReservationResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/penalty/{ticket_id} [put]
func (server *Server) cancelReservation(ctx *gin.Context) {
	var req cancelReservationRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)

	reservation, err := server.Queries.GetReservationDetails(ctx, req.TicketID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("reserved not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if reservation.UserID != authPayload.UserID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("don't have the peremision")))
		return
	}

	if reservation.Status != "RESERVED" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("ticket already is not-reserved")))
		return
	}

	err = server.Queries.CancelReservation(ctx, req.TicketID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	argReservation := db.UpdateReservationParams{
		Status: db.TicketStatusCANCELED,
		ID:     reservation.ID,
	}

	_, err = server.Queries.UpdateReservation(ctx, argReservation)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	penalty, err := server.Queries.GetTicketPenalties(ctx, reservation.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	departureTimeTicket := reservation.DepartureTime
	now := time.Now()

	if now.After(departureTimeTicket) {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("the cancellation time has passed")))
		return
	}

	timeRemaining := departureTimeTicket.Sub(now)

	oneHour := time.Hour

	var penaltyPercentage int32

	if timeRemaining <= oneHour {
		penaltyPercentage = penalty.AfterDay
	} else {
		penaltyPercentage = penalty.BeforDay
	}

	totalAmount := reservation.Amount * int64((100-penaltyPercentage)/100)

	argWallet := db.AddToUserWalletParams{
		Wallet: totalAmount,
		UserID: authPayload.UserID,
	}

	err = server.Queries.AddToUserWallet(ctx, argWallet)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	argPaymentStatus := db.UpdatePaymentStatusParams{
		Status: db.PaymentStatusFAILED,
		ID:     reservation.PaymentID,
	}

	_, err = server.Queries.UpdatePaymentStatus(ctx, argPaymentStatus)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	argPayment := db.UpdatePaymentAmountParams{
		Amount: reservation.Amount,
		ID:     reservation.PaymentID,
	}

	err = server.Queries.UpdatePaymentAmount(ctx, argPayment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	argChangeReservation := db.CreateChangeReservationParams{
		ReservationID: reservation.ID,
		AdminID:       1,
		UserID:        authPayload.UserID,
		FromStatus:    db.TicketStatusRESERVED,
		ToStatus:      db.TicketStatusCANCELED,
	}

	cngr, err := server.Queries.CreateChangeReservation(ctx, argChangeReservation)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	argTicket := db.UpdateTicketStatusParams{
		ID:     req.TicketID,
		Status: db.CheckReservationTicketStatusNOTRESERVED,
	}

	_, err = server.Queries.UpdateTicketStatus(ctx, argTicket)
	if err != nil {
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

	response := cancelReservationResponse{
		Message:           "CANCELED",
		AmountRefunded:    reservation.Amount,
		TicketID:          req.TicketID,
		ChangeReservation: cngr,
	}

	ctx.JSON(http.StatusOK, response)
}
