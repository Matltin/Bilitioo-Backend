package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

type SearchTicketsRequest struct {
	OriginCityID      int64  `form:"origin_city_id"`      // اختیاری
	DestinationCityID int64  `form:"destination_city_id"` // اختیاری
	DepartureDate     string `form:"departure_date"`      // تاریخ به صورت string گرفته میشه، بعداً تبدیل می‌کنیم
	VehicleType       string `form:"vehicle_type"`        // اختیاری
}

func (server *Server) searchTickets(ctx *gin.Context) {
	var req SearchTicketsRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// تبدیل تاریخ به time.Time اگر داده شده بود
	var departureDate time.Time
	if req.DepartureDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.DepartureDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid date format. Use YYYY-MM-DD")))
			return
		}
		departureDate = parsedDate
	}

	// params := db.SearchTicketsParams{
	// 	Column1: getInt64OrDefault(req.OriginCityID),
	// 	Column2: getInt64OrDefault(req.DestinationCityID),
	// 	Column3: getTimeOrDefault(departureDate),
	// 	Column4: getVehicleTypeOrDefault(req.VehicleType),
	// }

	params := db.SearchTicketsParams{
		Column1: req.OriginCityID,
		Column2: req.DestinationCityID,
		Column3: departureDate,
		Column4: db.VehicleType(req.VehicleType),
	}

	log.Println("\n\n", params, "\n\n")

	tickets, err := server.Queries.SearchTickets(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	ctx.JSON(http.StatusOK, tickets)
}

func getInt64OrDefault(p *int64) int64 {
	if p != nil {
		return *p
	}
	return 0
}

func getTimeOrDefault(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}

func getVehicleTypeOrDefault(vt *string) db.VehicleType {
	if vt != nil {
		return db.VehicleType(*vt)
	}
	return "" // یا مقدار default
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
