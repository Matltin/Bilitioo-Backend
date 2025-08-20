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
		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("authorization header is not provided")))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 || strings.ToLower(fields[0]) != authorizationTypeBearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("invalid authorization header")))
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
		payload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)

		user, err := server.Queries.GetUserByID(ctx, payload.UserID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if user.Role != "ADMIN" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(errors.New("access denied: admin privileges required")))
			return
		}

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
