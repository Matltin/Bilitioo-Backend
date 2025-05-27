package api

import (
	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config  util.Config
	router  *gin.Engine
	Queries *db.Queries
}

func NewServer(config util.Config, db *db.Queries) *Server {
	router := gin.Default()
	ser := &Server{
		config: config,
		Queries: db,
		router:  router,
	}

	ser.router.POST("/sign-in", ser.signUpUser)

	return ser
}

func (server *Server) Start(add string) {
	server.router.Run(add)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
