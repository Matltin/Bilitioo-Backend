package api

import (
	"database/sql"
	"net/http"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/token"
	"github.com/gin-gonic/gin"
)

func (server *Server) getReports(ctx *gin.Context) {
	reports, err := server.Queries.GetReports(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, reports)
}

type answerReportRequest struct {
	ID           int64  `json:"id"`
	ResponseText string `json:"response_text"`
}

func (server *Server) answerReport(ctx *gin.Context) {
	var req answerReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)

	arg := db.AnswerReportParams{
		ResponseText: req.ResponseText,
		ID:           req.ID,
		AdminID:      authPayload.UserID,
	}

	report, err := server.Queries.AnswerReport(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, report)
}

type updateTicketByReportRequest struct {
	ReservationID       int64  `json:"reserevation_id"`
	ToStatusReservation string `json:"to_status_reservation"`
}

func (server *Server) updateTicketByReport(ctx *gin.Context) {
	var req updateTicketByReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)

	status, err := server.Queries.GetReservationStatus(ctx, req.ReservationID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	reserve, err := server.Queries.GetReservationDetails(ctx, req.ReservationID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if req.ToStatusReservation == "CANCELED" || req.ToStatusReservation == "CANCELED-BY-TIME" {
		if status == "RESERVED" {
			ch := change{
				AdminID:             authPayload.UserID,
				Amount:              +reserve.Amount,
				ToStatusReservation: req.ToStatusReservation,
				Reserve:             reserve,
			}

			server.changeadd(ctx, ch)
		} else {
			ch := change{
				AdminID:             authPayload.UserID,
				Amount:              0,
				ToStatusReservation: req.ToStatusReservation,
				Reserve:             reserve,
			}

			server.chageWithOutAdd(ctx, ch)
		}
	} else if req.ToStatusReservation == "RESERVED" {
		ch := change{
			AdminID:             authPayload.UserID,
			Amount:              -reserve.Amount,
			ToStatusReservation: req.ToStatusReservation,
			Reserve:             reserve,
		}

		server.changeadd(ctx, ch)

	} else if req.ToStatusReservation == "RESERVING" {
		if status == "RESERVED" {
			ch := change{
				AdminID:             authPayload.UserID,
				Amount:              +reserve.Amount,
				ToStatusReservation: req.ToStatusReservation,
				Reserve:             reserve,
			}

			server.changeadd(ctx, ch)
		} else {
			ch := change{
				AdminID:             authPayload.UserID,
				Amount:              0,
				ToStatusReservation: req.ToStatusReservation,
				Reserve:             reserve,
			}

			server.chageWithOutAdd(ctx, ch)
		}
	}

	ctx.JSON(http.StatusOK, nil)
}

type change struct {
	AdminID             int64
	Amount              int64
	ToStatusReservation string
	Reserve             db.GetReservationDetailsRow
}

func (server *Server) changeadd(ctx *gin.Context, ch change) {
	argWallet := db.AddToUserWalletParams{
		Wallet: ch.Amount,
		UserID: ch.Reserve.UserID,
	}

	err := server.Queries.AddToUserWallet(ctx, argWallet)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	argPayment := db.UpdatePaymentAmountParams{
		Amount: ch.Amount,
		ID:     ch.Reserve.PaymentID,
	}

	err = server.Queries.UpdatePaymentAmount(ctx, argPayment)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.chageWithOutAdd(ctx, ch)
}

func (server *Server) chageWithOutAdd(ctx *gin.Context, ch change) {
	argReservation := db.UpdateReservationParams{
		Status: db.TicketStatus(ch.ToStatusReservation),
		ID:     ch.Reserve.ID,
	}
	status, err := server.Queries.GetReservationStatus(ctx, ch.Reserve.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	updatedReservation, err := server.Queries.UpdateReservation(ctx, argReservation)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	argChangeReservation := db.CreateChangeReservationParams{
		ReservationID: updatedReservation.ID,
		AdminID:       ch.AdminID,
		UserID:        updatedReservation.UserID,
		FromStatus:    status,
		ToStatus:      updatedReservation.Status,
	}

	_, err = server.Queries.CreateChangeReservation(ctx, argChangeReservation)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
}

type createReportRequest struct {
	RequestText   string `json:"request_text"`
	RequestType   string `json:"request_type"`
	ReservationID int64  `json:"reservation_id"`
}

func (server *Server) createReport(ctx *gin.Context) {
	var req createReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)

	arg := db.CreateReportParams{
		ReservationID: req.ReservationID,
		RequestType: db.RequestType(req.RequestType),
		RequestText: req.RequestText,
		UserID: authPayload.UserID,
	}
	
	report, err := server.Queries.CreateReport(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, report)
}
