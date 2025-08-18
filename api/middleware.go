package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Matltin/Bilitioo-Backend/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPyloadKey  = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPyloadKey, payload)
		ctx.Next()
	}
}

// adminMiddleware checks if the authenticated user has admin role
func (server *Server) adminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the auth payload (user must be authenticated first)
		authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)

		// Get user details from database
		user, err := server.Queries.GetUserByID(ctx, authPayload.UserID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// Check if user has admin role
		if user.Role != "ADMIN" {
			err := errors.New("access denied: admin privileges required")
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}

		// User is admin, continue to next handler
		ctx.Next()
	}
}

// Alternative: Combined auth + admin middleware for better performance
func (server *Server) adminAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// First check authentication
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := server.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Get user and check admin role
		user, err := server.Queries.GetUserByID(ctx, payload.UserID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if user.Role != "ADMIN" {
			err := errors.New("access denied: admin privileges required")
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}

		ctx.Set(authorizationPyloadKey, payload)
		ctx.Next()
	}
}
