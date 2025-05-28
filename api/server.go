package api

import (
	"database/sql"
	"log"
	"net/http"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/token"
	"github.com/Matltin/Bilitioo-Backend/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	router     *gin.Engine
	Queries    *db.Queries
	tokenMaker token.Maker
}

func NewServer(config util.Config, db *db.Queries) *Server {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetrickey)
	if err != nil {
		log.Fatal(err)
	}

	ser := &Server{
		config:     config,
		Queries:    db,
		tokenMaker: tokenMaker,
	}

	ser.setupRouter()

	return ser
}

func (ser *Server) setupRouter() {
	router := gin.Default()

	router.POST("/sign-in", ser.signUpUser)
	router.POST("/log-in", ser.logInUser)

	authRoutes := router.Group("/").Use(authMiddleware(ser.tokenMaker))

	authRoutes.PUT("/profile", ser.updateProfile)
	authRoutes.GET("/profile", ser.getUserProfile)
	authRoutes.GET("/city", ser.getCities)
	authRoutes.POST("/city", ser.searchTicketsByCities)
	authRoutes.GET("/ticket-detail/:ticket_id", ser.getTicketDetails)

	ser.router = router
}

func (server *Server) Start(add string) {
	server.router.Run(add)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
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
