package http

import (
	"net/http"

	"order-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase *usecase.OrderUsecase
}

func NewHandler(uc *usecase.OrderUsecase) *Handler {
	return &Handler{usecase: uc}
}

type OrderRequest struct {
	CustomerID string `json:"customer_id"`
	ItemName   string `json:"item_name"`
	Amount     int64  `json:"amount"`
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req OrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.usecase.CreateOrder(req.CustomerID, req.ItemName, req.Amount)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "payment service unavailable"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *Handler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	order, err := h.usecase.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *Handler) CancelOrder(c *gin.Context) {
	id := c.Param("id")

	err := h.usecase.CancelOrder(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order cancelled"})
}

func (h *Handler) GetStats(c *gin.Context) {
	status := c.Query("status")

	if status != "" {
		count, err := h.usecase.GetStatsByStatus(status)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error(400)": "invalid status"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"count": count})
		return
	}

	stats, _ := h.usecase.GetStats()
	c.JSON(http.StatusOK, stats)
}
