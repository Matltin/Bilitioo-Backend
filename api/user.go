package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/util"
	"github.com/Matltin/Bilitioo-Backend/worker"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
)

func isValidPhoneNumber(phone string) bool {
	matched, _ := regexp.MatchString(`^09\d{9}$`, phone)
	return matched
}

func isValidEmail(email string) bool {
	// Simple regex for email validation
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, email)
	return matched
}

type signUpUserRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password" binding:"required,min=8"`
}

// before add redis
func (server *Server) signUpUser(ctx *gin.Context) {
	var req signUpUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.Email != "" && !isValidEmail(req.Email) {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid email format")))
		return
	}

	if req.PhoneNumber != "" && !isValidPhoneNumber(req.PhoneNumber) {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid phone number format. It must start with 09 and be 11 digits long")))
		return
	}

	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var phoneVerify bool = false
	if req.PhoneNumber != "" {
		phoneVerify = true
	}

	if req.Email == "" && req.PhoneNumber == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("either email or phone_number must be provided")))
		return
	}

	arg := db.CreateUserParams{
		HashedPassword: hashedPassword,
		Email:          req.Email,
		PhoneNumber:    req.PhoneNumber,
		EmailVerified:  false,
		PhoneVerified:  phoneVerify,
	}

	user, err := server.Queries.CreateUser(ctx, arg)
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

	err = server.Queries.InitialProfile(ctx, user.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	payload := worker.PayloadSendVerfyEmail{
		Email: user.Email,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}

	server.distribution.DistributTaskSendVerifyEmail(ctx, &payload, opts...)
	ctx.JSON(http.StatusOK, nil)
}

type logInUserRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type logInUserResponse struct {
	User        db.GetUserRow `json:"user"`
	AccessToken string        `json:"access_token"`
}

func (server *Server) logInUser(ctx *gin.Context) {
	var req logInUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.Email == "" && req.PhoneNumber == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("either email or phone_number must be provided")))
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

	if user.EmailVerified == false {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("verify your email first")))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := logInUserResponse{
		User:        user,
		AccessToken: accessToken,
	}

	ctx.JSON(http.StatusOK, res)
}

// api/auth.go
func (server *Server) registerUserRedis(ctx *gin.Context) {
	var req signUpUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check if user exists in cache first (for rate limiting or duplicate checks)
	cacheKey := fmt.Sprintf("signup:attempt:%s:%s", req.Email, req.PhoneNumber)
	exists, err := server.redisClient.Exists(ctx, cacheKey)
	if err == nil && exists {
		ctx.JSON(http.StatusTooManyRequests, errorResponse(errors.New("please wait before trying again")))
		return
	}

	// Set a temporary key to prevent rapid duplicate signups
	server.redisClient.Set(ctx, cacheKey, "1", 20*time.Second)

	// Rest of your existing validation code...
	if req.Email != "" && !isValidEmail(req.Email) {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid email format")))
		return
	}

	// ... rest of your existing signup logic ...

	if req.PhoneNumber != "" && !isValidPhoneNumber(req.PhoneNumber) {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid phone number format. It must start with 09 and be 11 digits long")))
		return
	}

	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var phoneVerify bool = false
	if req.PhoneNumber != "" {
		phoneVerify = true
	}

	if req.Email == "" && req.PhoneNumber == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("either email or phone_number must be provided")))
		return
	}

	arg := db.CreateUserParams{
		HashedPassword: hashedPassword,
		Email:          req.Email,
		PhoneNumber:    req.PhoneNumber,
		EmailVerified:  false,
		PhoneVerified:  phoneVerify,
	}

	user, err := server.Queries.CreateUser(ctx, arg)
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

	err = server.Queries.InitialProfile(ctx, user.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	payload := worker.PayloadSendVerfyEmail{
		Email: user.Email,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}

	server.distribution.DistributTaskSendVerifyEmail(ctx, &payload, opts...)

	// After successful signup, you might want to cache the new user
	userCacheKey := fmt.Sprintf("user:%s:%s", user.Email, user.PhoneNumber)
	userJSON, err := json.Marshal(user)
	if err == nil {
		server.redisClient.Set(ctx, userCacheKey, string(userJSON), 5*time.Minute)
	}

	ctx.JSON(http.StatusOK, nil)

	// ... rest of your existing code ...
}

// api/auth.go
func (server *Server) loginUserRedis(ctx *gin.Context) {
	var req logInUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.Email == "" && req.PhoneNumber == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("either email or phone_number must be provided")))
		return
	}

	// Create a cache key
	cacheKey := fmt.Sprintf("user:%s:%s", req.Email, req.PhoneNumber)

	// Try to get from cache first
	cachedUser, err := server.redisClient.Get(ctx, cacheKey)

	if err == nil {
		var user db.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			// Verify password
			if err := util.CheckPassword(req.Password, user.HashedPassword); err != nil {
				ctx.JSON(http.StatusUnauthorized, errorResponse(err))
				return
			}

			accessToken, _, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			res := logInUserResponse{
				User: db.GetUserRow{
					ID:             user.ID,
					Email:          user.Email,
					PhoneNumber:    user.PhoneNumber,
					HashedPassword: user.HashedPassword,
				},
				AccessToken: accessToken,
			}

			ctx.JSON(http.StatusOK, res)
			return
		}
	}

	// Cache miss, proceed with database query
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

	// Verify password
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Cache the user data
	userJSON, err := json.Marshal(user)
	if err == nil {
		server.redisClient.Set(ctx, cacheKey, string(userJSON), 5*time.Minute)
	}

	accessToken, _, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := logInUserResponse{
		User:        user,
		AccessToken: accessToken,
	}

	ctx.JSON(http.StatusOK, res)
}
