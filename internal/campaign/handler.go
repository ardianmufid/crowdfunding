package campaign

import (
	"crowdfunding/internal/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type handler struct {
	svc service
	// svcUser user.Service
}

func NewHandler(svc service) handler {
	return handler{
		svc: svc,
	}
}

func (h handler) GetCampaigns(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Query("user_id"))
	if err != nil {
		response := helper.NewResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)

		ctx.JSON(http.StatusBadRequest, response)
		return
	}

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
