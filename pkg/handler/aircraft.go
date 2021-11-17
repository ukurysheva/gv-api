package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	gvapi "github.com/ukurysheva/gv-api"
)

func (h *Handler) createAircraft(c *gin.Context) {
	fmt.Println("createAircraft")
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input gvapi.Aircraft

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Wrong values")
		return
	}

	success := CheckAircraftValues(c, input)
	if !success {
		return
	}

	id, err := h.services.Aircraft.Create(userId, input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllAircrafts(c *gin.Context) {

	aircrafts, err := h.services.Aircraft.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllAircraftsResponse{
		Data: aircrafts,
	})
}

func (h *Handler) getAircraftById(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	aircraft, err := h.services.Aircraft.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, aircraft)
}

type getAllAircraftsResponse struct {
	Data []gvapi.Aircraft `json:"data"`
}

func CheckAircraftValues(c *gin.Context, input gvapi.Aircraft) bool {

	flagVals := make(map[string]bool)
	flagVals["Y"] = true
	flagVals["N"] = true

	if _, ok := flagVals[input.EconomyClass]; !ok {
		newErrorResponse(c, http.StatusBadRequest, "wrong value for EconomyClass")
		return false
	}
	if _, ok := flagVals[input.BusinessClass]; !ok {
		newErrorResponse(c, http.StatusBadRequest, "wrong value for BusinessClass")
		return false
	}

	if _, ok := flagVals[input.FirstClass]; !ok {
		newErrorResponse(c, http.StatusBadRequest, "wrong value for FirstClass")
		return false
	}

	return true

}
