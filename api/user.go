package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type signInUserRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func (server *Server) signInUser(ctx *gin.Context) {
	var req signInUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}
