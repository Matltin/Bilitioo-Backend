package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

const (
	ToAccount    = "myself"
	DefualtAdmin = 1
)

type createReservationRequest struct {
	Tickets []int64 `json:"tickets" binding:"required"`
}

func (server *Server) createReservation(ctx *gin.Context) {
	var req createReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if len(req.Tickets) == 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("there is no reservations")))
		return
	}

	var tickets []db.Ticket
	var amount int64 = 0
	for _, i := range req.Tickets {
		t, err := server.Queries.GetTicket(ctx, i)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		tickets = append(tickets, t)
		amount += t.Amount
	}

	authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)

	argPayment := db.CreatePaymentParams{
		FromAccount: authPayload.UserID,
		ToAccount:   ToAccount,
		Amount:      amount,
	}

	payment, err := server.Queries.CreatePayment(ctx, argPayment)
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

	var reservations []db.CreateReservationRow

	for _, t := range tickets {
		argReservation := db.CreateReservationParams{
			UserID:    authPayload.UserID,
			TicketID:  t.ID,
			PaymentID: payment.ID,
		}

		reservation, err := server.Queries.CreateReservation(ctx, argReservation)
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

		reservations = append(reservations, reservation)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"reservations": reservations,
		"payment":      payment,
	})
}

func (server *Server) getAllUserReservation(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)
	reservations, err := server.Queries.GetAllUserReservation(ctx, authPayload.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, reservations)
}
