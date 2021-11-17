package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	gvapi "github.com/ukurysheva/gv-api"
)

type getAllAirlineResponse struct {
	Data []gvapi.Airline `json:"data"`
}

func (h *Handler) createAirline(c *gin.Context) {
	fmt.Println("createAirline")
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input gvapi.Airline

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Wrong values")
		return
	}

	success := CheckAirlineValues(c, input)
	if !success {
		return
	}

	id, err := h.services.Airline.Create(userId, input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllAirlines(c *gin.Context) {

	airlines, err := h.services.Airline.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllAirlineResponse{
		Data: airlines,
	})
}

func (h *Handler) getAirlineById(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	airline, err := h.services.Airline.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, airline)
}

func CheckAirlineValues(c *gin.Context, input gvapi.Airline) bool {

	flagVals := make(map[string]bool)
	flagVals["Y"] = true
	flagVals["N"] = true
	fmt.Println(input)
	if _, ok := flagVals[input.Active]; !ok {
		newErrorResponse(c, http.StatusBadRequest, "wrong value for Active")
		return false
	}

	return true
}
