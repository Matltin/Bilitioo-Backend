package api

import (
	"log"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	db_redis "github.com/Matltin/Bilitioo-Backend/redis"
	"github.com/Matltin/Bilitioo-Backend/token"
	"github.com/Matltin/Bilitioo-Backend/util"
	"github.com/Matltin/Bilitioo-Backend/worker"
	"github.com/elastic/go-elasticsearch"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config        util.Config
	router        *gin.Engine
	Queries       *db.Queries
	tokenMaker    token.Maker
	distribution  worker.TaskDistributor
	redisClient   *db_redis.Client
	elasticClient *elasticsearch.Client
}

func NewServer(config util.Config, distributor worker.TaskDistributor, db *db.Queries, redis *db_redis.Client, es *elasticsearch.Client) *Server {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetrickey)
	if err != nil {
		log.Fatal(err)
	}

	ser := &Server{
		config:        config,
		elasticClient: es,
		Queries:       db,
		tokenMaker:    tokenMaker,
		distribution:  distributor,
		redisClient:   redis,
	}

	ser.setupRouter()

	return ser
}

func (ser *Server) setupRouter() {
	router := gin.Default()

	// Add CORS middleware - IMPORTANT for frontend integration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:8080",
		"http://127.0.0.1:8080",
		"http://localhost:3000", // For development
		"*",                     // For testing - remove in production
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Authorization",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers",
	}
	config.AllowCredentials = true

	router.Use(cors.New(config))

	swaggerUrl := ginSwagger.URL("http://localhost:3000/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerUrl))

	// Public routes (no authentication required)
	router.POST("/sign-in", ser.registerUserRedis)
	router.POST("/log-in", ser.loginUserRedis)
	router.GET("/verify-email", ser.verifyEmail)

	// Protected routes (authentication required)
	authRoutes := router.Group("/").Use(authMiddleware(ser.tokenMaker))
	{
		// User routes
		authRoutes.PUT("/profile", ser.updateProfile)
		authRoutes.GET("/profile", ser.getUserProfile)
		authRoutes.GET("/city", ser.getCities)
		authRoutes.POST("/city", ser.searchTicketsByCities)
		authRoutes.POST("/search-tickets", ser.searchTickets)
		authRoutes.GET("/ticket-detail/:ticket_id", ser.getTicketDetails)
		authRoutes.POST("/reservation", ser.createReservation)
		authRoutes.GET("/completedReservation", ser.getCompletedUserReservation)
		authRoutes.GET("/allReservation", ser.getAllUserReservation)
		authRoutes.POST("/payment", ser.payPayment)
		authRoutes.GET("/ticket-penalties/:ticket_id", ser.getTicketPenalties)
		authRoutes.GET("/penalty/:ticket_id", ser.getTicketPenalties)
		authRoutes.PUT("/penalty/:ticket_id", ser.cancelReservation)
		authRoutes.GET("/completed-tickets", ser.getAllUserCompletedTickets)
		authRoutes.GET("/notcompleted-tickets", ser.getAllUserNotCompletedTickets)
		authRoutes.POST("/report", ser.createReport) // Users can create reports
	}

	// Admin-only routes (requires both authentication and admin role)
	adminRoutes := router.Group("/admin").Use(authMiddleware(ser.tokenMaker))
	{
		adminRoutes.GET("/reports", ser.getReports)                  // Admin only
		adminRoutes.PUT("/reports/manage", ser.updateTicketByReport) // Admin only
		adminRoutes.PUT("/reports/answer", ser.answerReport)         // Admin only
		adminRoutes.GET("/tickets", ser.getAllTickets)               // Admin only

		adminRoutes.GET("/users/:userID/completed-tickets", ser.getCompletedTicketsForUserByAdmin)
		adminRoutes.GET("/users/:userID/notcompleted-tickets", ser.getNotCompletedTicketsForUserByAdmin)

	}

	// Alternative: Mixed routes where some endpoints need admin access
	mixedRoutes := router.Group("/").Use(authMiddleware(ser.tokenMaker))
	{
		// This endpoint can be accessed by both users and admins, but with different responses
		mixedRoutes.GET("/reports", ser.getReportsWithRoleCheck)
	}

	ser.router = router
}

func (server *Server) Start(add string) {
	server.router.Run(add)
}

func errorResponse(err error) gin.H {
	log.Println("error: ", err.Error())
	return gin.H{"error": err.Error()}
}
