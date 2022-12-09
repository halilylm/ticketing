package dto

type NewOrder struct {
	TicketID string `json:"ticket_id" validate:"required"`
}
