package transaction

import (
	"crowdfunding/internal/helper"
	"crowdfunding/internal/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	svc     service
	svcUser user.Service
}

func NewHandler(svc service, svcUser user.Service) handler {
	return handler{
		svc:     svc,
		svcUser: svcUser,
	}
}

func (h handler) GetCampaignTransactions(ctx *gin.Context) {

	var request GetCampaignTrasactionRequest

	err := ctx.ShouldBindUri(&request)
	if err != nil {
		response := helper.NewResponse("Failed to get campaign's transaction", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	currentUserID := ctx.GetInt("USER_ID")
	currentUser, err := h.svcUser.GetUserByID(ctx.Request.Context(), currentUserID)
	if err != nil {
		response := helper.NewResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	request.User = currentUser

	transactions, err := h.svc.GetTransactionByCampaignID(ctx.Request.Context(), request)
	if err != nil {
		response := helper.NewResponse("Failed to get campaign's transaction", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.NewResponse("Campaign's transactions", http.StatusOK, "success", NewMapperCampaignTransactionsResponse(transactions))
	ctx.JSON(http.StatusOK, response)
}

func (h handler) GetUserTransactions(ctx *gin.Context) {

	userID := ctx.GetInt("USER_ID")

	transactions, err := h.svc.GetTransactionByUserID(ctx.Request.Context(), userID)
	if err != nil {
		response := helper.NewResponse("Failed to get user's transaction", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.NewResponse("User's transactions", http.StatusOK, "success", NewMapperUserTransactionsResponse(transactions))
	ctx.JSON(http.StatusOK, response)
}

func (h handler) CreateTransaction(ctx *gin.Context) {

	var request CreateTransactionRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMassage := gin.H{"errors": errors}

		response := helper.NewResponse("Failed to create transaction", http.StatusUnprocessableEntity, "errors", errorMassage)

		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUserID := ctx.GetInt("USER_ID")
	currentUser, err := h.svcUser.GetUserByID(ctx.Request.Context(), currentUserID)
	if err != nil {
		response := helper.NewResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	request.User = currentUser

	transaction, err := h.svc.CreateTransaction(ctx.Request.Context(), request)
	if err != nil {
		response := helper.NewResponse("Failed to create transaction", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.NewResponse("Success to create transactions", http.StatusOK, "success", NewMapperTransactionResponse(transaction))
	ctx.JSON(http.StatusOK, response)
}
