package api

import (
	"database/sql"
	"log"
	"net/http"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

type getTicketDetailsRequest struct {
	TicketID int64 `uri:"ticket_id" binding:"required"`
}

func (server *Server) getTicketDetails(ctx *gin.Context) {
	var req getTicketDetailsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log.Println("\n\n\n", req.TicketID, "\n\n\n")

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
