package api

import (
	"database/sql"
	"errors"
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

	totalAmount := reservation.Amount * int64((100 - penaltyPercentage)/100)

	argWallet := db.AddToUserWalletParams{
		Wallet: totalAmount,
		UserID: authPayload.UserID,
	}

	err = server.Queries.AddToUserWallet(ctx, argWallet)
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

	response := cancelReservationResponse{
		Message:           "CANCELED",
		AmountRefunded:    reservation.Amount,
		TicketID:          req.TicketID,
		ChangeReservation: cngr,
	}

	ctx.JSON(http.StatusOK, response)
}
