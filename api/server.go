package api

import (
	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router  *gin.Engine
	Queries *db.Queries
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

func errorResponse(err error) gin.H {
  return gin.H{"error": err.Error()}
}
