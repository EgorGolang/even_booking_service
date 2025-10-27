package handler

import (
	"event_booking_service/internal/models"
	"event_booking_service/internal/service"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Борнирование
func (h *Handler) CreateBooking(c *gin.Context) {
	var req models.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: fmt.Sprintf("binding json error: %v", err),
		})
		return
	}
	booking, err := h.service.CreateBooking(req)
	if err != nil {
		errorMsg := err.Error()
		statusCode := http.StatusUnprocessableEntity
		switch {
		case strings.Contains(errorMsg, "not enough tickets"):
			statusCode = http.StatusConflict
		case strings.Contains(errorMsg, "not found"):
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, booking)
		return
	}
	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Booking Created",
		Data:    booking,
	})
}

/*func (h *Handler) GetBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "id is not a number",
		})
		return
	}
	BookingWithEvent, err := h.service.GetBooking(c.Request.Context(), id)
	if err != nil {
		fmt.Printf("Filed to get booking", err)
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Booking Not Found",
		})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Booking Success",
		Data:    BookingWithEvent,
	})
}*/

func (h *Handler) GetUserBooking(c *gin.Context) {
	userIDstr := c.Param("user_id")
	if userIDstr == "" {
		fmt.Printf("userID  is empty\n")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "userID is empty",
		})
		return
	}
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		fmt.Printf("userID is invalid\n")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "userID is invalid",
		})
		return
	}
	bookings, err := h.service.GetUserBookings(c.Request.Context(), userID)
	if err != nil {
		fmt.Printf("userID is invalid\n %s", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "userID is invalid",
		})
	}
	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Booking Success",
		Data:    bookings,
	})
}

func (h *Handler) CancelBooking(c *gin.Context) {
	var req models.CancelBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "bad request",
		})
		return
	}
	err := h.service.CancelBooking(c.Request.Context(), req.BookingID, req.UserID)
	if err != nil {
		fmt.Printf("cancel booking error: %v\n", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "cancel booking error",
		})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Booking Cancelled",
	})
}
