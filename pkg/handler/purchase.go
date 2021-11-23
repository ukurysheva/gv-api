package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	gvapi "github.com/ukurysheva/gv-api"
)

// type getAllAirlineResponse struct {
// 	Data []gvapi.Purchase `json:"data"`
// }

func (h *Handler) createPurchase(c *gin.Context) {
	fmt.Println("createAirline")
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input gvapi.Purchase

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Wrong values")
		return
	}

	success := CheckPurchaseValues(c, input)
	if !success {
		return
	}

	id, err := h.services.Purchase.Create(userId, input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// func (h *Handler) getAllAirlines(c *gin.Context) {

// 	airlines, err := h.services.Airline.GetAll()
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	c.JSON(http.StatusOK, getAllAirlineResponse{
// 		Data: airlines,
// 	})
// }

func (h *Handler) getPurchaseById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	purchase, err := h.services.Purchase.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, purchase)
}

func CheckPurchaseValues(c *gin.Context, input gvapi.Purchase) bool {

	flagVals := make(map[string]bool)
	flagVals["Y"] = true
	flagVals["N"] = true

	class := make(map[string]bool)
	class["economy"] = true
	class["pr_economy"] = true
	class["business"] = true
	class["first"] = true
	fmt.Println(input)
	if _, ok := flagVals[input.Food]; !ok {
		newErrorResponse(c, http.StatusBadRequest, "wrong value for Food")
		return false
	}

	if _, ok := class[input.Class]; !ok {
		newErrorResponse(c, http.StatusBadRequest, "wrong value for Class")
		return false
	}

	return true
}
