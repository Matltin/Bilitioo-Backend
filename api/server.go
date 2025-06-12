package api

import (
	"log"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	db_redis "github.com/Matltin/Bilitioo-Backend/redis"
	"github.com/Matltin/Bilitioo-Backend/token"
	"github.com/Matltin/Bilitioo-Backend/util"
	"github.com/Matltin/Bilitioo-Backend/worker"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config       util.Config
	router       *gin.Engine
	Queries      *db.Queries
	tokenMaker   token.Maker
	distribution worker.TaskDistributor
	redisClient  *db_redis.Client
}

func NewServer(config util.Config, distributor worker.TaskDistributor, db *db.Queries, redis *db_redis.Client) *Server {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetrickey)
	if err != nil {
		log.Fatal(err)
	}

	ser := &Server{
		config:       config,
		Queries:      db,
		tokenMaker:   tokenMaker,
		distribution: distributor,
		redisClient:  redis,
	}

	ser.setupRouter()

	return ser
}

func (ser *Server) setupRouter() {
	router := gin.Default()

	router.POST("/sign-in", ser.registerUserRedis) //1
	router.POST("/log-in", ser.loginUserRedis)  //2
	router.GET("/verify-email", ser.verifyEmail)

	authRoutes := router.Group("/").Use(authMiddleware(ser.tokenMaker))

	authRoutes.PUT("/profile", ser.updateProfile) //3
	authRoutes.GET("/profile", ser.getUserProfile)  //3
	authRoutes.GET("/city", ser.getCities)  //4
	authRoutes.GET("/search-tickets", ser.searchTickets) //5
	authRoutes.POST("/city", ser.searchTicketsByCities)
	authRoutes.GET("/ticket-detail/:ticket_id", ser.getTicketDetails)  //6
	authRoutes.POST("/reservation", ser.createReservation)  //7
	authRoutes.GET("/completedReservation", ser.getCompletedUserReservation)
	authRoutes.GET("/allReservation", ser.getAllUserReservation)
	authRoutes.POST("/payment", ser.payPayment)  //8


	authRoutes.GET("/ticket-penalties/:ticket_id", ser.getTicketPenalties)  //9
	authRoutes.GET("/penalty/:ticket_id", ser.getTicketPenalties) //9
	authRoutes.PUT("/penalty/:ticket_id", ser.cancelReservation)  //9, 12

	authRoutes.GET("/report", ser.getReports) //10
	authRoutes.PUT("/report", ser.answerReport) //13
	authRoutes.POST("/report", ser.createReport) //13
	authRoutes.PUT("/manage-report", ser.updateTicketByReport)  //10

	authRoutes.GET("/completed-tickets", ser.getAllUserCompletedTickets) //11
	authRoutes.GET("/notcompleted-tickets", ser.getAllUserNotCompletedTickets) //11
	authRoutes.GET("/tickets", ser.getAllTickets) //11

	ser.router = router
}

func (server *Server) Start(add string) {
	server.router.Run(add)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
