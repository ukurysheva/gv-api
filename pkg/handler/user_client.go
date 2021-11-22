package handler

// func (h *Handler) createClient(c *gin.Context) {
// 	fmt.Println("createCountry")
// 	userId, err := getUserId(c)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	var input gvapi.Country

// 	if err := c.BindJSON(&input); err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	id, err := h.services.Country.Create(userId, input)

// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	c.JSON(http.StatusOK, map[string]interface{}{
// 		"id": id,
// 	})
// }

// func (h *Handler) getAllUserClients(c *gin.Context) {

// 	countries, err := h.services.Country.GetAll()
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	c.JSON(http.StatusOK, getAllUserClientsResponse{
// 		Data: countries,
// 	})
// }

// func (h *Handler) getUserClientById(c *gin.Context) {

// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
// 		return
// 	}

// 	country, err := h.services.Country.GetById(id)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	c.JSON(http.StatusOK, country)
// }

// type getAllUserClientsResponse struct {
// 	Data []gvapi.Client `json:"data"`
// }
