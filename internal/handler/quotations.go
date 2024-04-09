package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type quotationResponse struct {
	Code       string    `json:"code"`
	Rate       float32   `json:"rate"`
	UpdateTime time.Time `json:"update_time"`
}

// @Summary Получить котировку по идентификатору
// @Description Возвращает котировку по идентификатору обновления
// @Produce json
// @Param id query string false "Идентификатор обновления"
// @Tags quotations
// @Success 200 {object} quotationResponse
// @Router /quotations/get [get]
func (h *Handler) getQuotation(c *gin.Context) {
	updateIDString := c.Query("id")
	if updateIDString == "" {
		newErrorResponse(c, http.StatusBadRequest, "error: provide id to find")
		return
	}

	updateID, err := strconv.Atoi(updateIDString)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "error: incorrect form of id")
		return
	}

	q, err := h.service.GetQuotation(updateID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, quotationResponse{
		Code:       fmt.Sprintf("%s/%s", q.CodeFrom, q.CodeTo),
		Rate:       q.Rate,
		UpdateTime: q.UpdateTime,
	})
}

// @Summary Обновить котировку
// @Description Создаёт новый запрос на обновление котировки
// @Produce json
// @Param from query string false "Базовая валюта"
// @Param to query string false "Котируемая валюта"
// @Tags quotations
// @Success 200 {object} map[string]any
// @Router /quotations/update [get]
func (h *Handler) updateQuotation(c *gin.Context) {
	codeFrom := strings.ToUpper(c.Query("from"))
	codeTo := strings.ToUpper(c.Query("to"))

	if codeFrom == "" || codeTo == "" {
		newErrorResponse(c, http.StatusBadRequest, "error: provide both parameters (from, to)")
		return
	}

	updateID, err := h.service.InsertQuotation(codeFrom, codeTo)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"update_id": updateID,
	})
}

// @Summary Получить последнее значение котировки
// @Description Возвращает последнее валидное значение котировки
// @Produce json
// @Param from query string false "Базовая валюта"
// @Param to query string false "Котируемая валюта"
// @Tags quotations
// @Success 200 {object} quotationResponse
// @Router /quotations/latest [get]
func (h *Handler) getLatestQuotation(c *gin.Context) {
	codeFrom := strings.ToUpper(c.Query("from"))
	codeTo := strings.ToUpper(c.Query("to"))

	if codeFrom == "" || codeTo == "" {
		newErrorResponse(c, http.StatusBadRequest, "error: provide both parameters (from, to)")
		return
	}

	q, err := h.service.GetLatestQuotation(codeFrom, codeTo)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, quotationResponse{
		Code:       fmt.Sprintf("%s/%s", q.CodeFrom, q.CodeTo),
		Rate:       q.Rate,
		UpdateTime: q.UpdateTime,
	})
}
