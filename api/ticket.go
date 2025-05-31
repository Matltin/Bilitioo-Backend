package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/token"
	"github.com/gin-gonic/gin"
)

// type SearchTicketsRequest struct {
// 	OriginCityID      *int64          `form:"origin_city_id"`
// 	DestinationCityID *int64          `form:"destination_city_id"`
// 	DepartureDate     *string         `form:"departure_date"`
// 	VehicleType       *db.VehicleType `form:"vehicle_type"`
// }

// func (server *Server) searchTickets(ctx *gin.Context) {
// 	var req SearchTicketsRequest

// 	if err := ctx.ShouldBindQuery(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	// تبدیل تاریخ به time.Time اگر داده شده بود
// 	var departureDate *time.Time
// 	if req.DepartureDate != nil {
// 		parsedDate, err := time.Parse("2006-01-02", *req.DepartureDate)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid date format. Use YYYY-MM-DD")))
// 			return
// 		}
// 		departureDate = &parsedDate
// 	}

// 	// مقداردهی پیش‌فرض برای پارامترهای NULL
// 	var originCityID, destinationCityID sql.NullInt64

// 	if req.OriginCityID != nil {
// 		originCityID = sql.NullInt64{Int64: *req.OriginCityID, Valid: true}
// 	}
// 	if req.DestinationCityID != nil {
// 		destinationCityID = sql.NullInt64{Int64: *req.DestinationCityID, Valid: true}
// 	}

// 	params := db.SearchTicketsParams{
// 		Column1: originCityID.Int64,
// 		Column2: destinationCityID.Int64,
// 		Column3: time.Time{},
// 		Column4: *req.VehicleType,
// 	}

// 	if departureDate != nil {
// 		params.Column3 = *departureDate // dereference the pointer
// 	}

// 	tickets, err := server.Queries.SearchTickets(ctx, params)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, tickets)
// }

// func getInt64OrDefault(p *int64) int64 {
// 	if p != nil {
// 		return *p
// 	}
// 	return 0
// }

// func getTimeOrDefault(t *time.Time) time.Time {
// 	if t != nil {
// 		return *t
// 	}
// 	return time.Time{}
// }

// func getVehicleTypeOrDefault(vt *string) db.VehicleType {
// 	if vt != nil {
// 		return db.VehicleType(*vt)
// 	}
// 	return "" // یا مقدار default
// }

type searchTicketsRequest struct {
	OriginCityID      int64  `json:"origin_city_id" binding:"required"`
	DestinationCityID int64  `json:"destination_city_id" binding:"required"`
	DepartureDate     string `json:"departure_date" binding:"required"`
	VehicleType       string `json:"vehicle_type" binding:"required,oneof=BUS TRAIN AIRPLANE"`
}

func (server *Server) searchTickets(ctx *gin.Context) {
	var req searchTicketsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Parse departure date
	departureDate, err := time.Parse("2006-01-02", req.DepartureDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid date format, expected YYYY-MM-DD")))
		return
	}

	// Calculate date range
	startOfDay := time.Date(departureDate.Year(), departureDate.Month(), departureDate.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	arg := db.SearchTicketsParams{
		OriginCityID:      req.OriginCityID,
		DestinationCityID: req.DestinationCityID,
		DepartureTime:     startOfDay,
		DepartureTime_2:   endOfDay,
		Column5:           db.VehicleType(req.VehicleType),
	}

	tickets, err := server.Queries.SearchTickets(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tickets)
}

type getTicketDetailsRequest struct {
	TicketID int64 `uri:"ticket_id" binding:"required"`
}

func (server *Server) getTicketDetails(ctx *gin.Context) {
	var req getTicketDetailsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ticket, err := server.Queries.GetTicketDetails(ctx, req.TicketID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := gin.H{
		"origin":        ticket.Origin,
		"destination":   ticket.Destination,
		"departureTime": ticket.DepartureTime,
		"arrivalTime":   ticket.ArrivalTime,
		"amount":        ticket.Amount,
		"capacity":      ticket.Capacity,
		"vehicle_type":  ticket.VehicleType,
		"feature":       ticket.Feature,
		"status":        ticket.Status,
	}

	switch ticket.VehicleType {
	case db.VehicleTypeBUS:
		response["VIP"] = ticket.VIP.Bool
		response["bed_chair"] = ticket.BedChair.Bool
	case db.VehicleTypeTRAIN:
		response["have_compartment"] = ticket.HaveCompartment.Bool
		response["rank"] = ticket.Rank.Int32
	case db.VehicleTypeAIRPLANE:
		response["flight_class"] = string(ticket.FlightClass.FlightClass)
		response["airplane_name"] = ticket.AirplaneName.String
	}

	ctx.JSON(http.StatusOK, response)
}

func (server *Server) getAllUserCompletedTickets(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)
	tickets, err := server.Queries.GetAllUserCompletedTickets(ctx, authPayload.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("reserved not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tickets)
}

func (server *Server) getAllUserNotCompletedTickets(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)
	tickets, err := server.Queries.GetAllUserNotCompletedTickets(ctx, authPayload.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("reserved not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tickets)
}

func (server *Server) getAllTickets(ctx *gin.Context) {
	tickets, err := server.Queries.GetAllTickets(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("reserved not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, tickets)
}
