package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type signUpUserRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password" binding:"required,min=8"`
}

func (server *Server) signUpUser(ctx *gin.Context) {
	var req signUpUserRequest
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

type logInUserRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type logInUserResponse struct {
	User db.GetUserRow		`json:"user"`
	AccessToken string	`json:"access_token"`
}

func (server *Server) logInUser(ctx *gin.Context) {
	var req logInUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.Email == "" && req.PhoneNumber == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("Either email or phone_number must be provided.")))
		return
	}

	arg := db.GetUserParams{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}

	user, err := server.Queries.GetUser(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}
	
	accessToken, _, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := logInUserResponse{
		User: user,
		AccessToken: accessToken,
	}

	ctx.JSON(http.StatusOK, res)
}

