package repository

import (
	"context"
	"database/sql"
	"event_booking_service/internal/models"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Новое бронирование
func (r *Repository) CreateBooking(ctx context.Context, booking *models.Booking) error {
	query := `
        INSERT INTO bookings (user_id, event_id, tickets, total_price, status)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRowContext(ctx, query,
		booking.UserID,
		booking.EventID,
		booking.Tickets,
		booking.TotalPrice,
		booking.Status,
	).Scan(&booking.ID, &booking.CreatedAt, &booking.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}

	return nil
}

// Бронирование оп ID
func (r *Repository) GetBookingByID(ctx context.Context, id int) (*models.Booking, error) {
	query := `
        SELECT id, user_id, event_id, tickets, total_price, status, created_at, updated_at
        FROM bookings WHERE id = $1
    `

	var booking models.Booking
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.EventID,
		&booking.Tickets,
		&booking.TotalPrice,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("booking not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	return &booking, nil
}

// Все бронирования пользователя
func (r *Repository) GetUserBookings(ctx context.Context, userID int) ([]models.Booking, error) {
	query := `
        SELECT id, user_id, event_id, tickets, total_price, status, created_at, updated_at
        FROM bookings 
        WHERE user_id = $1 
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user bookings: %w", err)
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.EventID,
			&booking.Tickets,
			&booking.TotalPrice,
			&booking.Status,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

// Отмена бронирования
func (r *Repository) CancelBooking(ctx context.Context, bookingID, userID int) error {
	query := `
        UPDATE bookings 
        SET status = 'cancelled' 
        WHERE id = $1 AND user_id = $2 AND status = 'confirmed'
    `

	result, err := r.db.ExecContext(ctx, query, bookingID, userID)
	if err != nil {
		return fmt.Errorf("failed to cancel booking: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("booking not found or already cancelled")
	}

	return nil
}

// Колличество доступных билетов
func (r *Repository) UpdateEventTickets(ctx context.Context, eventID, tickets int) error {
	query := `
		UPDATE events 
        SET available_tickets = available_tickets - $1 
        WHERE id = $2 AND available_tickets >= $1`
	result, err := r.db.ExecContext(ctx, query, tickets, eventID)
	if err != nil {
		return fmt.Errorf("failed to update events: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("event not found or already cancelled")
	}
	return nil
}

// Обновление бронирования
func (r *Repository) UpdateBooking(ctx context.Context, booking *models.Booking) error {
	query := `
        UPDATE bookings 
        SET tickets = $1, total_price = $2, status = $3
        WHERE id = $4 AND user_id = $5
    `
	result, err := r.db.ExecContext(ctx, query,
		booking.Tickets,
		booking.TotalPrice,
		booking.Status,
		booking.ID,
		booking.UserID)

	if err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("booking not found or already cancelled")
	}
	return nil
}

/*// Event cache methods

// SaveEventCache сохраняет информацию о мероприятии в кэш
func (r *Repository) SaveEventCache(ctx context.Context, event *models.Event) error {
	query := `
        INSERT INTO event_cache (id, title, date, total_tickets, available_tickets, price)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (id)
        DO UPDATE SET
            title = $2,
            date = $3,
            total_tickets = $4,
            available_tickets = $5,
            price = $6,
            last_updated = CURRENT_TIMESTAMP
    `

	_, err := r.db.ExecContext(ctx, query,
		event.ID,
		event.Title,
		event.Date,
		event.TotalTickets,
		event.AvailableTickets,
		event.Price,
	)

	return err
}

// GetEventCache возвращает информацию о мероприятии из кэша
func (r *Repository) GetEventCache(ctx context.Context, eventID int) (*models.Event, error) {
	query := `
        SELECT id, title, date, total_tickets, available_tickets, price
        FROM event_cache WHERE id = $1
    `

	var event models.Event
	err := r.db.QueryRowContext(ctx, query, eventID).Scan(
		&event.ID,
		&event.Title,
		&event.Date,
		&event.TotalTickets,
		&event.AvailableTickets,
		&event.Price,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("event not found in cache")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get event cache: %w", err)
	}

	return &event, nil
}
*/
