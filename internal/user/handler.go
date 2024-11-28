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

		response := helper.NewResponse("Register account failed", http.StatusBadRequest, "error", errorMessage)

		ctx.JSON(http.StatusBadRequest, response)
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

	registerUserResponse := NewRegisterUserResponse(newUser, "initokenbro")

	response := helper.NewResponse("Account has been registered", http.StatusCreated, "success", registerUserResponse)

	ctx.JSON(http.StatusOK, response)
}
