package payments

import "github.com/google/uuid"

type PayInvoiceRequest struct {
	InvoiceID uuid.UUID `json:"invoice_id" form:"invoice_id" validate:"required,uuid"`
}
