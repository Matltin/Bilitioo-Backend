package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

// verifyEmail godoc
//
//	@Summary		Verify user email
//	@Description	Confirms email ownership by validating a verification code. Requires authentication if user is logged in.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/verify-email [post]
func (server *Server) verifyEmail(ctx *gin.Context) {
	// Get query parameters
	idParam := ctx.Query("id")
	secretCode := ctx.Query("secret_code")

	if idParam == "" || secretCode == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("missing id or secret code")))
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid id format")))
		return
	}

	// Try to update the verify_emails table
	verifyEmail, err := server.Queries.UpdateVerifyEmail(ctx, db.UpdateVerifyEmailParams{
		ID:         id,
		SecretCode: secretCode,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid or expired verification link")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Mark user's email as verified
	err = server.Queries.UpdateUserEmailVerified(ctx, db.UpdateUserEmailVerifiedParams{
		ID:            verifyEmail.UserID,
		EmailVerified: true,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.Queries.GetUserByID(ctx, verifyEmail.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.invalidateUserCache(ctx, user.ID, user.Email, user.PhoneNumber)

	ctx.JSON(http.StatusOK, gin.H{"message": "email successfully verified"})
}
