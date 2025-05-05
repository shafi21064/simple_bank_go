package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/shafi21064/simplebank/db/sqlc"
	"github.com/shafi21064/simplebank/util"
)

// create account
type createUserRequest struct {
	UserName string `json:"user_name" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	UserName  string
	FullName  string
	Email     string
	CreatedAt pgtype.Timestamptz
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hassedPassword, err := util.HassedPassword(req.Password)
	if err != nil {
		util.CheckError("Failed to hassed password", err)
		return
	}
	arg := db.CreateUsersParams{
		UserName:       req.UserName,
		HassedPassword: hassedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUsers(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := createUserResponse{
		UserName:  user.UserName,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	ctx.JSON(http.StatusCreated, response)

}
