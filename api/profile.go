package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/token"
	"github.com/Matltin/Bilitioo-Backend/util"
	"github.com/gin-gonic/gin"
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
		var err error
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

	ctx.JSON(http.StatusOK, gin.H{
		"profile": profile,
		"user":    user,
	})
}

func (server *Server) getUserProfile(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)

	profile, err := server.Queries.GetUserProfile(ctx, authPayload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, profile)
}
