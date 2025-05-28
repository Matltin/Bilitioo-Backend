package api

import (
	"database/sql"
	"net/http"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/token"
	"github.com/gin-gonic/gin"
)

type updateProfileRequest struct {
  PicDir       string `json:"pic_dir"`
  FirstName    string `json:"first_name"`
  LastName     string `json:"last_name"`
  CityID       int64  `json:"city_id"`
  NationalCode string `json:"national_code"`
  PhoneNumber  string `json:"phone_number"`
  Email        string `json:"email"`
}


func (server *Server) updateProfile(ctx *gin.Context) {
  var req updateProfileRequest
  if err := ctx.ShouldBindJSON(&req); err != nil {
    ctx.JSON(http.StatusBadRequest, errorResponse(err))
    return
  }

  authPayload := ctx.MustGet(authorizationPyloadKey).(*token.Payload)

  profileArgs := db.UpdateProfileParams{
    UserID: authPayload.UserID,
    PicDir: sql.NullString{
      String: req.PicDir,
      Valid:  true,
    },
    FirstName: sql.NullString{
      String: req.FirstName,
      Valid:  true,
    },
    LastName: sql.NullString{
      String: req.LastName,
      Valid:  true,
    },
    CityID: sql.NullInt64{
      Int64: req.CityID,
      Valid: true,
    },
    NationalCode: sql.NullString{
      String: req.NationalCode,
      Valid:  true,
    },
  }

  profile, err := server.Queries.UpdateProfile(ctx, profileArgs)
  if err != nil {
    ctx.JSON(http.StatusInternalServerError, errorResponse(err))
    return
  }

  userArgs := db.UpdateUserContactParams{
    ID: authPayload.UserID,
    Email: sql.NullString{
      String: req.Email,
      Valid: true,
    },
    PhoneNumber: sql.NullString{
      String: req.PhoneNumber,
      Valid: true,
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