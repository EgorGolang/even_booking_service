package models

import (
	"time"
)

type Booking struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	EventID    int       `json:"event_id"`
	Tickets    int       `json:"tickets"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateBookingRequest struct {
	UserID  int `json:"user_id" binding:"required"`
	EventID int `json:"event_id" binding:"required"`
	Tickets int `json:"tickets" binding:"required,min=1,max=10"`
}

type CancelBookingRequest struct {
	BookingID int `json:"booking_id" binding:"required"`
	UserID    int `json:"user_id" binding:"required"`
}

type BookingWithEvent struct {
	Booking *Booking `json:"booking"`
	Event   *Event   `json:"event"`
}

type Event struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	Date             time.Time `json:"date"`
	TotalTickets     int       `json:"total_tickets"`
	AvailableTickets int       `json:"available_tickets"`
	Price            int       `json:"price"`
}
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}
