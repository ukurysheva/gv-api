package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	gvapi "github.com/ukurysheva/gv-api"
)

type getAllFlightsResponse struct {
	Data []gvapi.Flight `json:"data"`
}

func (h *Handler) createFlight(c *gin.Context) {
	fmt.Println("createFlight")
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input gvapi.Flight

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Wrong values")
		fmt.Println(err.Error())
		return
	}

	success := CheckFlightValues(c, input)
	if !success {
		return
	}

	id, err := h.services.Flight.Create(userId, input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllFlights(c *gin.Context) {

	flights, err := h.services.Flight.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllFlightsResponse{
		Data: flights,
	})
}

func (h *Handler) getFlightById(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	flight, err := h.services.Flight.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, flight)
}

func CheckFlightValues(c *gin.Context, input gvapi.Flight) bool {

	flagVals := make(map[string]bool)
	flagVals["Y"] = true
	flagVals["N"] = true

	if _, ok := flagVals[input.Wifi]; !ok {
		newErrorResponse(c, http.StatusBadRequest, "wrong value for Wifi")
		return false
	}
	if _, ok := flagVals[input.Usb]; !ok {
		newErrorResponse(c, http.StatusBadRequest, "wrong value for Usb")
		return false
	}

	if _, ok := flagVals[input.Food]; !ok {
		newErrorResponse(c, http.StatusBadRequest, "wrong value for Food")
		return false
	}

	return true

}
