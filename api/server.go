package api

import (
	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/token"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router  *gin.Engine
	Queries *db.Queries
	tokenMaker token.Maker
}

func NewServer(db *db.Queries) *Server {
	router := gin.Default()
	ser := &Server{
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
