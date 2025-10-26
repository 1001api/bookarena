package bookings

import (
	"time"

	"github.com/google/uuid"
)

type CreateBookingRequest struct {
	FieldID   uuid.UUID `json:"field_id" form:"field_id" validate:"required,uuid"`
	StartTime time.Time `json:"start_time" form:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" form:"end_time" validate:"required"`
}
