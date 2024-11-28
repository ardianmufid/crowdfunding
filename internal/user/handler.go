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

func (h handler) checkEmailAvailability(ctx *gin.Context) {
	var input CheckEmailRequest

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.NewResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)

		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.svc.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server Error"}
		response := helper.NewResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{"is_available": isEmailAvailable}

	metaMessage := "Email has been registered"
	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.NewResponse(metaMessage, http.StatusOK, "success", data)

	ctx.JSON(http.StatusOK, response)

}
