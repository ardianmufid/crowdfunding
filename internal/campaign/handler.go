package campaign

import (
	"crowdfunding/internal/helper"
	"crowdfunding/internal/user"
	"net/http"
	"strconv"

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

func (h handler) GetCampaigns(ctx *gin.Context) {
	userID, _ := strconv.Atoi(ctx.Query("user_id"))

	campaigns, err := h.svc.GetAllCampaign(userID)
	if err != nil {
		response := helper.NewResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)

		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.NewResponse("List of Campaigns", http.StatusOK, "success", NewMapperCampaignsResponse(campaigns))
	ctx.JSON(http.StatusOK, response)
}

func (h handler) GetCampaign(ctx *gin.Context) {

	var request CampaignDetailRequest

	err := ctx.ShouldBindUri(&request)
	if err != nil {
		response := helper.NewResponse("Failed to get detail campaign", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	campaign, err := h.svc.GetCampaignByID(request)
	if err != nil {
		response := helper.NewResponse("Failed to get detail campaign", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.NewResponse("Campaign detail", http.StatusOK, "success", NewMapperCampaignDetailResponse(campaign))
	ctx.JSON(http.StatusOK, response)
}

func (h handler) CreateCampaign(ctx *gin.Context) {

	var request CreateCampaignRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.NewResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUserID := ctx.GetInt("USER_ID")
	currentUser, err := h.svcUser.GetUserByID(currentUserID)

	if err != nil {
		response := helper.NewResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	request.User = currentUser

	newCampaign, err := h.svc.CreateCampaign(request)

	if err != nil {
		response := helper.NewResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.NewResponse("Success to crate campaign", http.StatusCreated, "success", NewMapperCampaignResponse(newCampaign))
	ctx.JSON(http.StatusCreated, response)

}
