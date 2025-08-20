package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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

// searchTickets godoc
//
//	@Summary		Search tickets
//	@Description	Search tickets by origin, destination, date, and vehicle type
//	@Tags			tickets
//	@Accept			json
//	@Produce		json
//	@Param			request	body		searchTicketsRequest	true	"Search tickets request"
//	@Success		200		{array}		db.Ticket
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/tickets/search [post]
func (server *Server) searchTickets(ctx *gin.Context) {
	var req searchTicketsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cacheKey := fmt.Sprintf("search:%d:%d:%s:%s", req.OriginCityID, req.DestinationCityID, req.DepartureDate, req.VehicleType)

	// Try Redis
	cached, err := server.redisClient.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var tickets []db.Ticket
		if err := json.Unmarshal([]byte(cached), &tickets); err == nil {
			ctx.JSON(http.StatusOK, tickets)
			return
		}
	}

	// Parse and search DB
	departureDate, err := time.Parse("2006-01-02", req.DepartureDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid date format, expected YYYY-MM-DD")))
		return
	}

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

	// Save to Redis
	jsonData, _ := json.Marshal(tickets)
	server.redisClient.Set(ctx, cacheKey, jsonData, 5*time.Minute)

	ctx.JSON(http.StatusOK, tickets)
}

type getTicketDetailsRequest struct {
	TicketID int64 `uri:"ticket_id" binding:"required"`
}

// getTicketDetails godoc
//
//	@Summary		Get ticket details
//	@Description	Retrieve detailed info for a ticket by ID
//	@Tags			tickets
//	@Accept			json
//	@Produce		json
//	@Param			ticket_id	path		int64	true	"Ticket ID"
//	@Success		200			{object}	map[string]interface{}
//	@Failure		400			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/tickets/{ticket_id} [get]
func (server *Server) getTicketDetails(ctx *gin.Context) {
	var req getTicketDetailsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cacheKey := fmt.Sprintf("ticket_details:%d", req.TicketID)

	// Check Redis
	cached, err := server.redisClient.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var cachedResponse map[string]interface{}
		if err := json.Unmarshal([]byte(cached), &cachedResponse); err == nil {
			ctx.JSON(http.StatusOK, cachedResponse)
			return
		}
	}

	// DB Query fallback
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

	// Store in Redis
	jsonData, _ := json.Marshal(response)
	server.redisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute)

	ctx.JSON(http.StatusOK, response)
}

// getAllUserCompletedTickets godoc
//
//	@Summary		Get completed tickets for user
//	@Description	Get all tickets the authenticated user has completed
//	@Tags			tickets
//	@Produce		json
//	@Success		200	{array}		db.Ticket
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/tickets/user/completed [get]
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

// getAllUserNotCompletedTickets godoc
//
//	@Summary		Get not completed tickets for user
//	@Description	Get all tickets the authenticated user has reserved but not completed
//	@Tags			tickets
//	@Produce		json
//	@Success		200	{array}		db.Ticket
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/tickets/user/not-completed [get]
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

// getAllTickets godoc
//
//	@Summary		Get all tickets
//	@Description	Retrieve all tickets from the system
//	@Tags			tickets
//	@Produce		json
//	@Success		200	{array}		db.Ticket
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/tickets [get]
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

func (server *Server) searchTicketsElastic(ctx *gin.Context) {
	var req searchTicketsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{"term": map[string]interface{}{"origin_city_id": req.OriginCityID}},
					map[string]interface{}{"term": map[string]interface{}{"destination_city_id": req.DestinationCityID}},
					map[string]interface{}{"term": map[string]interface{}{"vehicle_type.keyword": req.VehicleType}},
					map[string]interface{}{
						"range": map[string]interface{}{
							"departure_time": map[string]interface{}{
								"gte": req.DepartureDate,
							},
						},
					},
				},
			},
		},
	}

	// (postgres).ticket  1.api 2.   -> [database elastic ]  -> searchTicketsElastic (result)

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res, err := server.elasticClient.Search(
		server.elasticClient.Search.WithContext(ctx),
		server.elasticClient.Search.WithIndex("tickets"),
		server.elasticClient.Search.WithBody(&buf),
		server.elasticClient.Search.WithTrackTotalHits(true),
		server.elasticClient.Search.WithPretty(),
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, r)
}