package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	gvapi "github.com/ukurysheva/gv-api"
)

type getAllArportResponse struct {
	Data []gvapi.Airport `json:"data"`
}

func (h *Handler) createAirport(c *gin.Context) {
	fmt.Println("createAirport")
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input gvapi.Airport

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Wrong values")
		return
	}

	success := CheckAirportValues(c, input)
	if !success {
		return
	}

	id, err := h.services.Airport.Create(userId, input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllAirports(c *gin.Context) {

	airports, err := h.services.Airport.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllArportResponse{
		Data: airports,
	})
}

func (h *Handler) getAirportById(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	airport, err := h.services.Airport.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, airport)
}

func CheckAirportValues(c *gin.Context, input gvapi.Airport) bool {

	flagVals := make(map[string]bool)
	flagVals["Y"] = true
	flagVals["N"] = true

	typesVals := make(map[string]bool)
	typesVals["balloonport"] = true
	typesVals["closed"] = true
	typesVals["heliport"] = true
	typesVals["large_airport"] = true
	typesVals["medium_airport"] = true
	typesVals["small_airport"] = true

	if _, ok := flagVals[input.Visa]; !ok {
		newErrorResponse(c, http.StatusBadRequest, "wrong value for Visa")
		return false
	}
	if _, ok := flagVals[input.CovidTest]; !ok {
		newErrorResponse(c, http.StatusBadRequest, "wrong value for CovidTest")
		return false
	}

	if _, ok := flagVals[input.Quarantine]; !ok {
		newErrorResponse(c, http.StatusBadRequest, "wrong value for Quarantine")
		return false
	}

	if _, ok := flagVals[input.LockDown]; !ok {
		newErrorResponse(c, http.StatusBadRequest, "wrong value for LockDown")
		return false
	}

	if _, ok := typesVals[input.Type]; !ok {
		newErrorResponse(c, http.StatusBadRequest, "wrong value for Type")
		return false
	}

	return true

}
