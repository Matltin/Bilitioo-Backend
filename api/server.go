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

	router := gin.Default()
	ser := &Server{
		config:  config,
		Queries: db,
		router:  router,
		tokenMaker: tokenMaker,
	}


	ser.router.POST("/sign-in", ser.signUpUser)
	ser.router.POST("/log-in", ser.logInUser)

	return ser
}

func (server *Server) Start(add string) {
	server.router.Run(add)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
