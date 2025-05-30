package api

import (
	"log"

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
	authRoutes.POST("/reservation", ser.createReservation)
	authRoutes.POST("/payment", ser.payPayment)

	authRoutes.GET("/search-tickets", ser.searchTickets)

	ser.router = router
}

func (server *Server) Start(add string) {
	server.router.Run(add)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
