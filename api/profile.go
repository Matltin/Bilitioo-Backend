package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/token"
	"github.com/Matltin/Bilitioo-Backend/util"
	"github.com/Matltin/Bilitioo-Backend/worker"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

type updateProfileRequest struct {
	PicDir       string `json:"pic_dir"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	CityID       int64  `json:"city_id"`
	NationalCode string `json:"national_code"`
	PhoneNumber  string `json:"phone_number"`
	Password     string `json:"password"`
	Email        string `json:"email"`
}

func (server *Server) updateProfile(ctx *gin.Context) {
	var req updateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)

	// Get current user data to check what's changing
	currentUser, err := server.Queries.GetUserByID(ctx, authPayload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Validate inputs
	if req.Email != "" && !isValidEmail(req.Email) {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid email format")))
		return
	}

	if req.PhoneNumber != "" && !isValidPhoneNumber(req.PhoneNumber) {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid phone number format. It must start with 09 and be 11 digits long")))
		return
	}

	profileArgs := db.UpdateProfileParams{
		UserID: authPayload.UserID,
		PicDir: sql.NullString{
			String: req.PicDir,
			Valid:  req.PicDir != "",
		},
		FirstName: sql.NullString{
			String: req.FirstName,
			Valid:  req.FirstName != "",
		},
		LastName: sql.NullString{
			String: req.LastName,
			Valid:  req.LastName != "",
		},
		CityID: sql.NullInt64{
			Int64: req.CityID,
			Valid: req.CityID > 0,
		},
		NationalCode: sql.NullString{
			String: req.NationalCode,
			Valid:  req.NationalCode != "",
		},
	}

	profile, err := server.Queries.UpdateProfile(ctx, profileArgs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var newPass string = ""
	if req.Password != "" {
		newPass, err = util.HashedPassword(req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	userArgs := db.UpdateUserContactParams{
		ID: authPayload.UserID,
		Email: sql.NullString{
			String: req.Email,
			Valid:  req.Email != "",
		},
		PhoneNumber: sql.NullString{
			String: req.PhoneNumber,
			Valid:  req.PhoneNumber != "",
		},
		HashedPassword: sql.NullString{
			String: newPass,
			Valid:  req.Password != "",
		},
	}

	user, err := server.Queries.UpdateUserContact(ctx, userArgs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if req.Email != "" {
		payload := worker.PayloadSendVerfyEmail{
			Email: user.Email,
		}

		opts := []asynq.Option{
			asynq.MaxRetry(10),
			asynq.ProcessIn(10 * time.Second),
			asynq.Queue(worker.QueueCritical),
		}

		server.distribution.DistributTaskSendVerifyEmail(ctx, &payload, opts...)

		arg := db.UpdateUserEmailVerifiedParams{
			ID:            profile.UserID,
			EmailVerified: false,
		}
		err := server.Queries.UpdateUserEmailVerified(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	// Invalidate Redis cache for this user
	server.invalidateUserCache(ctx, authPayload.UserID, currentUser.Email, currentUser.PhoneNumber)

	err = server.cacheUserData(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"profile": profile,
		"user":    user,
	})
}

func (server *Server) invalidateUserCache(ctx *gin.Context, userID int64, oldEmail, oldPhoneNumber string) {
	// Delete cache entries that might exist for this user
	cacheKeys := []string{
		fmt.Sprintf("user:%d", userID),                      // ID-based key
		fmt.Sprintf("user:%s:%s", oldEmail, oldPhoneNumber), // Old email key
		fmt.Sprintf("profile:%d", userID),                   // Profile key
	}

	for _, key := range cacheKeys {
		if err := server.redisClient.Delete(ctx, key); err != nil {
			// Log the error but don't fail the request
			log.Printf("Failed to delete cache key %s: %v", key, err)
		}
	}
}

func (server *Server) cacheUserData(ctx *gin.Context, user db.User) error {
	// Create multiple cache keys for different lookup methods
	profile, err := server.Queries.GetUserProfile(ctx, user.ID)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("profile:%d", profile.UserID)

	keys := []string{
		fmt.Sprintf("user:%d", user.ID),                         // ID-based key
		fmt.Sprintf("user:%s:%s", user.Email, user.PhoneNumber), // Email-based key
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	for _, key := range keys {
		if err := server.redisClient.Set(ctx, key, userJSON, 24*time.Hour); err != nil {
			return err
		}
	}
	profileJSON, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	err = server.redisClient.Set(ctx, cacheKey, profileJSON, 24*time.Hour)
	if err != nil {
		return err
	}

	return nil
}

// func (server *Server) updateProfile(ctx *gin.Context) {
// 	var req updateProfileRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)

// 	if req.Email != "" && !isValidEmail(req.Email) {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid email format")))
// 		return
// 	}

// 	if req.PhoneNumber != "" && !isValidPhoneNumber(req.PhoneNumber) {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid phone number format. It must start with 09 and be 11 digits long")))
// 		return
// 	}

// 	profileArgs := db.UpdateProfileParams{
// 		UserID: authPayload.UserID,
// 		PicDir: sql.NullString{
// 			String: req.PicDir,
// 			Valid:  req.PicDir != "",
// 		},
// 		FirstName: sql.NullString{
// 			String: req.FirstName,
// 			Valid:  req.FirstName != "",
// 		},
// 		LastName: sql.NullString{
// 			String: req.LastName,
// 			Valid:  req.LastName != "",
// 		},
// 		CityID: sql.NullInt64{
// 			Int64: req.CityID,
// 			Valid: req.CityID > 0,
// 		},
// 		NationalCode: sql.NullString{
// 			String: req.NationalCode,
// 			Valid:  req.NationalCode != "",
// 		},
// 	}

// 	profile, err := server.Queries.UpdateProfile(ctx, profileArgs)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	var newPass string = ""

// 	if req.Password != "" {
// 		var err error
// 		newPass, err = util.HashedPassword(req.Password)
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 			return
// 		}
// 	}

// 	userArgs := db.UpdateUserContactParams{
// 		ID: authPayload.UserID,
// 		Email: sql.NullString{
// 			String: req.Email,
// 			Valid:  req.Email != "",
// 		},
// 		PhoneNumber: sql.NullString{
// 			String: req.PhoneNumber,
// 			Valid:  req.PhoneNumber != "",
// 		},
// 		HashedPassword: sql.NullString{
// 			String: newPass,
// 			Valid:  req.Password != "",
// 		},
// 	}

// 	user, err := server.Queries.UpdateUserContact(ctx, userArgs)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"profile": profile,
// 		"user":    user,
// 	})
// }

// func (server *Server) getUserProfile(ctx *gin.Context) {
// 	authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)

// 	profile, err := server.Queries.GetUserProfile(ctx, authPayload.UserID)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, profile)
// }

func (server *Server) getUserProfile(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)

	// Try to get from cache first
	cacheKey := fmt.Sprintf("profile:%d", authPayload.UserID)
	cachedProfile, err := server.redisClient.Get(ctx, cacheKey)
	if err == nil {
		var profile db.Profile
		if err := json.Unmarshal([]byte(cachedProfile), &profile); err == nil {
			ctx.JSON(http.StatusOK, profile)
			return
		}
	}

	fmt.Println("hello")

	// Cache miss - get from database
	profile, err := server.Queries.GetUserProfile(ctx, authPayload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Cache the profile data
	profileJSON, err := json.Marshal(profile)
	if err == nil {
		if err := server.redisClient.Set(ctx, cacheKey, profileJSON, 24*time.Hour); err != nil {
			log.Printf("Failed to cache profile: %v", err)
		}
	}

	ctx.JSON(http.StatusOK, profile)
}
