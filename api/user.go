package api

import (
	"errors"
	"net/http"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type signInUserRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type logInUserRequest struct {
	Email string	`json:"email"`
	PhoneNumber string	`json:"phone_number"`
	Password string	`json:"password"`
}

func (server *Server) signInUser(ctx *gin.Context) {
	var req signInUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var emailVerify bool = false
	if req.Email != "" {
		emailVerify = true
	}

	var phoneVerify bool = false
	if req.PhoneNumber != "" {
		phoneVerify = true
	}

	if req.Email == "" && req.PhoneNumber == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("Either email or phone_number must be provided.")))
		return
	}

	arg := db.CreateUserParams{
		HashedPassword: hashedPassword,
		Email:          req.Email,
		PhoneNumber:    req.PhoneNumber,
		EmailVerified:  emailVerify,
		PhoneVerified:  phoneVerify,
	}

	_, err = server.Queries.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, nil)
}


