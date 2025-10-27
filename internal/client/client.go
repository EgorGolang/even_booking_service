package client

import (
	"bytes"
	"context"
	"encoding/json"
	"event_booking_service/internal/models"
	"fmt"
	"io"
	"net/http"
)

type EventService struct {
	baseURL string
	client  *http.Client
}

func NewEventService(baseURL string, client *http.Client) *EventService {
	return &EventService{
		baseURL: baseURL,
		client:  client,
	}
}

// Информация о меропричятии
func (r *EventService) GetEvent(ctx context.Context, id int) (*models.Booking, error) {
	resp, err := r.client.Get(fmt.Sprintf("%s/%d/event", r.baseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var event models.Booking
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

func (c *EventService) ReserveBooking(eventID, tickets int) error {
	requestBody := map[string]int{"tickets": tickets}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed json marshaling: %v", err)
	}
	resp, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%d/reserve", c.baseURL, eventID), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *EventService) ReleaseTickets(eventID, tickets int) error {
	requestBody := map[string]int{"tickets": tickets}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%d/reserve", c.baseURL, eventID), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed json marshaling: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to release tickets: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to release tickets: %d", resp.StatusCode)
	}
	return nil
}
