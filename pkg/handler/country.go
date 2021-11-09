package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	gvapi "github.com/ukurysheva/gv-api"
)

func (h *Handler) createCountry(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input gvapi.Country

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Country.Create(userId, input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllCountries(c *gin.Context) {

	countries, err := h.services.Country.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllCountriresResponse{
		Data: countries,
	})
}

type getAllCountriresResponse struct {
	Data []gvapi.Country `json:"data"`
}
