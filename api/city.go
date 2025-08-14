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

const userActivityID = "userActivityID"

type CityResponse struct {
	Province string `json:"province"`
	County   string `json:"county"`
}

func (server *Server) getCities(ctx *gin.Context) {
	cities, err := server.Queries.GetCities(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cities)
}

type searchTicketsByCitiesRequest struct {
	OriginCityID      int64  `json:"origin_city_id" binding:"required"`
	DestinationCityID int64  `json:"destination_city_id" binding:"required"`
	VehicleType       string `json:"vehicle_type" binding:"required,oneof=BUS AIRPLANE TRAIN"`
}

func (server *Server) searchTicketsByCities(ctx *gin.Context) {
	var req searchTicketsByCitiesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	argCities := db.SearchTicketsByCitiesParams{
		OriginCityID:      req.OriginCityID,
		DestinationCityID: req.DestinationCityID,
		VehicleType:       db.VehicleType(req.VehicleType),
	}

	tickets, err := server.Queries.SearchTicketsByCities(ctx, argCities)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(tickets) == 0 {
		ctx.JSON(http.StatusNotFound, errorResponse(errors.New("there is no ticket with this route")))
		return
	}

	authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)

	argUserActivity := db.CreateUserActivityParams{
		UserID:      authPayload.UserID,
		RouteID:     tickets[0].RouteID,
		VehicleType: db.VehicleType(req.VehicleType),
	}

	userActivity, err := server.Queries.CreateUserActivity(ctx, argUserActivity)
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

	ctx.Set(userActivityID, userActivity.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"tickets":      tickets,
		userActivityID: userActivity.ID,
	})
}
