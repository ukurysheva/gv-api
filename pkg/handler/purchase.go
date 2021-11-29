package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	gvapi "github.com/ukurysheva/gv-api"
)

func (h *Handler) createPurchase(c *gin.Context) {
	fmt.Println("createAirline")
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input gvapi.Purchase

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Ошибка: Некорректные данные")
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

func (h *Handler) getPurchaseById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Ошибка: Некорректное значение id билета")
		return
	}

	purchase, err := h.services.Purchase.GetById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, purchase)
}

func (h *Handler) getUserPurchases(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	purchase, err := h.services.Purchase.GetByUserId(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, purchase)
}

func (h *Handler) getUserBasket(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	purchase, err := h.services.Purchase.GetBasketByUserId(userId)
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
		newErrorResponse(c, http.StatusBadRequest, "Ошибка: Неправильно задан параметр 'Питание включено'")
		return false
	}

	if _, ok := class[input.Class]; !ok {
		newErrorResponse(c, http.StatusBadRequest, "Ошибка: Неправильно задано имя класса билета")
		return false
	}

	return true
}

func (h *Handler) payPurchase(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	purchaseId, err := strconv.Atoi(c.Param("id"))
	if purchaseId <= 0 {
		fmt.Println("BindJSON UpdateUserInput error")
		newErrorResponse(c, http.StatusBadRequest, "Ошибка: Некорректный id билета.")
		return
	}

	purchase, err := h.services.Purchase.GetById(purchaseId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Ошибка: Не удалось найти забронированный билет.")
		return
	}
	if purchase.Payed == 1 {
		newErrorResponse(c, http.StatusBadRequest, "Ошибка: Билет уже оплачен. Приятного полета!")
		return
	}
	if mins, err := strconv.Atoi(purchase.BookTimeLeft); err != nil || mins <= 0 {
		fmt.Println(err)
		newErrorResponse(c, http.StatusBadRequest, "Ошибка: Невозможно купить билет, так как срок действия брони истек")
		return
	}

	userProfile, err := h.services.User.GetProfile(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = CheckUserProfileData(userProfile)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input gvapi.PurchasePayInput
	if err := c.BindJSON(&input); err != nil {
		fmt.Println("BindJSON UpdateUserInput error")
		newErrorResponse(c, http.StatusBadRequest, "Ошибка: Некорректные данные для оплаты")
		return
	}

	err = CheckUserPayData(input, userProfile)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.PurchaseId = purchaseId

	if err := h.services.Purchase.Update(input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func CheckUserProfileData(userProfile gvapi.User) error {
	if userProfile.FirstName == "" || userProfile.LastName == "" || userProfile.MiddleName == "" || userProfile.BirthDate == "" ||
		userProfile.LivingAddress == "" || userProfile.PassportAddress == "" || userProfile.PassportNumber == "" ||
		userProfile.PassportSeries == "" {

		err := errors.New("Ошибка: Не заполнены персональные данные пользователя")

		return err
	}

	return nil
}

func CheckUserPayData(purchaseDetails gvapi.PurchasePayInput, userProfile gvapi.User) error {
	if purchaseDetails.PayMethod == "card" &&
		(userProfile.CardNumber == "" || userProfile.CardExpDate == "" || userProfile.CardIndividual == "") {
		err := errors.New("Ошибка: Не заполнены данные карты пользователя, необходимые для оплаты")
		return err
	}
	return nil
}
