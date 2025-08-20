package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

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

// getCompletedTicketsForUserByAdmin godoc
// @Summary      Get a user's completed tickets (Admin)
// @Description  Admin gets all completed tickets for a specific user by their ID.
// @Tags         Admin
// @Produce      json
// @Param        userID path int true "User ID"
// @Success      200 {array} db.GetAllUserCompletedTicketsForAdminRow
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /admin/users/{userID}/completed-tickets [get]
func (server *Server) getCompletedTicketsForUserByAdmin(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("userID"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid user ID")))
		return
	}

	tickets, err := server.Queries.GetAllUserCompletedTicketsForAdmin(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("no completed tickets found for this user")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tickets)
}

// getNotCompletedTicketsForUserByAdmin godoc
// @Summary      Get a user's pending tickets (Admin)
// @Description  Admin gets all pending (not completed) tickets for a specific user by their ID.
// @Tags         Admin
// @Produce      json
// @Param        userID path int true "User ID"
// @Success      200 {array} db.GetAllUserNotCompletedTicketsForAdminRow
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /admin/users/{userID}/notcompleted-tickets [get]
func (server *Server) getNotCompletedTicketsForUserByAdmin(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("userID"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid user ID")))
		return
	}

	tickets, err := server.Queries.GetAllUserNotCompletedTicketsForAdmin(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("no pending tickets found for this user")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tickets)
}
