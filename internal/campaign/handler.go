package campaign

import (
	"crowdfunding/internal/helper"
	"crowdfunding/internal/user"
	"fmt"
	"log"
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

func (h handler) CreateCampaign(ctx *gin.Context) {

	var request CreateCampaignRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMassage := gin.H{"errors": errors}

		response := helper.NewResponse("Failed to create campaign", http.StatusUnprocessableEntity, "errors", errorMassage)

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

	newCampaign, err := h.svc.CreateCampaign(ctx.Request.Context(), request)
	if err != nil {
		response := helper.NewResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.NewResponse("Success to create campaign", http.StatusCreated, "success", NewMapperCampaignResponse(newCampaign))
	ctx.JSON(http.StatusCreated, response)

}

func (h handler) GetCampaigns(ctx *gin.Context) {
	userID, _ := strconv.Atoi(ctx.Query("user_id"))

	campaigns, err := h.svc.GetAllCampaign(ctx.Request.Context(), userID)
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

	campaign, err := h.svc.GetCampaignByID(ctx.Request.Context(), request)
	if err != nil {
		response := helper.NewResponse("Failed to get detail campaign", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.NewResponse("Campaign detail", http.StatusOK, "success", NewMapperCampaignDetailResponse(campaign))
	ctx.JSON(http.StatusOK, response)
}

func (h handler) UpdateCampaign(ctx *gin.Context) {

	var requestID CampaignDetailRequest
	var requestData CreateCampaignRequest

	err := ctx.ShouldBindUri(&requestID)
	if err != nil {
		log.Printf("handler ShouldBindUri error : %v", err)
		response := helper.NewResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err = ctx.ShouldBindJSON(&requestData)
	if err != nil {
		log.Printf("handler ShouldBindJSON error : %v", err)
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.NewResponse("Failed to update campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUserID := ctx.GetInt("USER_ID")
	currentUser, err := h.svcUser.GetUserByID(ctx.Request.Context(), currentUserID)
	if err != nil {
		log.Printf("handler currentUser error : %v", err)
		response := helper.NewResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	requestData.User = currentUser

	updatedCampaign, err := h.svc.UpdateCampaign(ctx.Request.Context(), requestID, requestData)
	if err != nil {
		log.Printf("handler updatedCampaign error : %v", err)
		response := helper.NewResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.NewResponse("Success to update campaign", http.StatusOK, "success", NewMapperCampaignResponse(updatedCampaign))
	ctx.JSON(http.StatusOK, response)

}

func (h handler) UploadImage(ctx *gin.Context) {
	var request CreateCampaignImageRequest

	err := ctx.ShouldBind(&request)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.NewResponse("Failed to upload campaign image", http.StatusUnprocessableEntity, "error", errorMessage)
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

	file, err := ctx.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.NewResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	path := fmt.Sprintf("images/%d-%s", currentUserID, file.Filename)

	err = ctx.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.NewResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.svc.SaveCampaignImage(ctx.Request.Context(), request, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.NewResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.NewResponse("Campaign image succesfully uploaded", http.StatusCreated, "success", data)
	ctx.JSON(http.StatusCreated, response)
}
