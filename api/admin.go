package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/token"
	"github.com/gin-gonic/gin"
)

// Helper function to check if user is admin
func (server *Server) isUserAdmin(ctx *gin.Context, userID int64) (bool, error) {
	user, err := server.Queries.GetUserByID(ctx, userID)
	if err != nil {
		return false, err
	}
	return user.Role == "ADMIN", nil
}

// getReportsWithRoleCheck - Example of role-based access within handler
func (server *Server) getReportsWithRoleCheck(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)

	// Check if user is admin
	isAdmin, err := server.isUserAdmin(ctx, authPayload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !isAdmin {
		// Regular users can only see their own reports
		reports, err := server.Queries.GetUserReport(ctx, authPayload.UserID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(errors.New("no reports found")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, reports)
		return
	}

	// Admin can see all reports
	reports, err := server.Queries.GetReports(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, reports)
}

// updateAnswerReport with admin check
func (server *Server) answerReport(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)

	// Check if user is admin
	isAdmin, err := server.isUserAdmin(ctx, authPayload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !isAdmin {
		ctx.JSON(http.StatusForbidden, errorResponse(errors.New("access denied: admin privileges required")))
		return
	}

	// Continue with original answerReport logic
	var req answerReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.AnswerReportParams{
		ResponseText: req.ResponseText,
		ID:           req.ID,
		AdminID:      authPayload.UserID,
	}

	report, err := server.Queries.AnswerReport(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, report)
}
