package user

import (
	"crowdfunding/internal/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	svc service
}

func newHandler(svc service) handler {
	return handler{
		svc: svc,
	}
}

func (h handler) register(ctx *gin.Context) {

	var request RegisterUserRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.NewResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage)

		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.svc.RegisterUser(request)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.NewResponse("Register account failed", http.StatusBadRequest, "error", errorMessage)

		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	registerUserResponse := NewMapperUserResponse(newUser, "initokenbro")

	response := helper.NewResponse("Account has been registered", http.StatusCreated, "success", registerUserResponse)

	ctx.JSON(http.StatusOK, response)
}

func (h handler) login(ctx *gin.Context) {
	var request LoginUserRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.NewResponse("login failed", http.StatusUnprocessableEntity, "error", errorMessage)

		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.svc.LoginUser(request)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.NewResponse("login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loginUserResponse := NewMapperUserResponse(loggedinUser, "initokenbro")

	response := helper.NewResponse("Succesfully loggedin", http.StatusOK, "success", loginUserResponse)

	ctx.JSON(http.StatusOK, response)
}
