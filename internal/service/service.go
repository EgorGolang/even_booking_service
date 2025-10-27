package service

import (
	"context"
	"event_booking_service/internal/client"
	"event_booking_service/internal/models"
	"event_booking_service/internal/repository"
	"fmt"
)

type Service struct {
	repo        *repository.Repository
	eventClient *client.EventService
}

func NewService(repo *repository.Repository, eventClient *client.EventService) *Service {
	return &Service{
		repo:        repo,
		eventClient: eventClient,
	}
}

func (s *Service) CreateBooking(req models.CreateBookingRequest) (*models.Booking, error) {
	if err := s.eventClient.ReserveBooking(req.EventID, req.Tickets); err != nil {
		return nil, fmt.Errorf("failed to reserve tickets: %v", err)
	}
	booking := &models.Booking{
		UserID:     req.UserID,
		EventID:    req.EventID,
		Tickets:    req.Tickets,
		Status:     "confirmed",
	}
	return booking, nil
}

func (s *Service) GetUserBookings(ctx context.Context, userID int) ([]*models.BookingWithEvent, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user id %v", userID)
	}
	bookings, err := s.repo.GetUserBookings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings: %v", err)
	}
	var result []*models.BookingWithEvent
	for _, booking := range bookings {
		result = append(result, &models.BookingWithEvent{
			Booking: &booking,
			//Event:   &event,
		})
	}
	return result, nil
}

func (s *Service) CancelBooking(ctx context.Context, userID, bookingID int) error {
	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("failed to get booking: %v", err)
	}
	if booking.Status == "cancelled" {
		return fmt.Errorf("booking %v is already cancelled", bookingID)
	}
	if err := s.eventClient.ReleaseTickets(booking.EventID, booking.Tickets); err != nil {
		return fmt.Errorf("failed to release tickets: %v", err)
	}
	if err := s.repo.CancelBooking(ctx, bookingID, userID); err != nil {
		return fmt.Errorf("failed to cancel booking: %v", err)
	}
	return nil
}
