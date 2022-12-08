package http

import (
	"net/http"

	"github.com/halilylm/gommon/middlewares"
	"github.com/halilylm/gommon/rest"
	"github.com/halilylm/gommon/utils"
	"github.com/halilylm/ticketing/tickets/domain"
	"github.com/halilylm/ticketing/tickets/ticket/usecase"
	"github.com/labstack/echo/v4"
)

type ticketHandler struct {
	ticketUC usecase.Ticket
}

// NewTicketHandler handler for auth
func NewTicketHandler(g *echo.Group, ticketUC usecase.Ticket) {
	handler := &ticketHandler{ticketUC: ticketUC}

	// jwt middleware
	g.Use(middlewares.CurrentUser("jwt"))

	g.POST("/", handler.NewTicket)
	g.PUT("/:id", handler.UpdateTicket)
	g.GET("/:id", handler.ShowTicket)
	g.GET("/", handler.AvailableTickets)
}

// NewTicket creates a new ticket
func (t *ticketHandler) NewTicket(c echo.Context) error {
	var ticket domain.Ticket

	// bind request body to ticket
	if err := c.Bind(&ticket); err != nil {
		return c.JSON(rest.ErrorResponse(rest.NewBadRequestError(err.Error())))
	}

	// validate the struct
	if err := utils.ValidateStruct(&ticket); err != nil {
		return c.JSON(rest.ErrorResponse(rest.NewValidationErrors(err)))
	}

	// fill user id read from cookie
	user := middlewares.UserFromContext(c)
	ticket.UserID = user.ID

	// call the usecase
	createdTicket, err := t.ticketUC.NewTicket(c.Request().Context(), &ticket)
	if err != nil {
		// errors returning from usecase layer will be rest errors
		// so err can be used directly
		return c.JSON(rest.ErrorResponse(err))
	}

	return c.JSON(http.StatusCreated, createdTicket)
}

// UpdateTicket updates an existing ticket
func (t *ticketHandler) UpdateTicket(c echo.Context) error {
	var ticket domain.Ticket

	// bind request body to ticket
	if err := c.Bind(&ticket); err != nil {
		return c.JSON(rest.ErrorResponse(rest.NewBadRequestError(err.Error())))
	}

	// validate the struct
	if err := utils.ValidateStruct(&ticket); err != nil {
		return c.JSON(rest.ErrorResponse(rest.NewValidationErrors(err)))
	}

	// fill the ticket id
	ticket.ID = c.Param("id")

	// call the usecase
	updatedTicket, err := t.ticketUC.UpdateTicket(c.Request().Context(), &ticket)
	if err != nil {
		// errors returning from usecase layer will be rest errors
		// so err can be used directly
		return c.JSON(rest.ErrorResponse(err))
	}

	return c.JSON(http.StatusOK, updatedTicket)
}

func (t *ticketHandler) ShowTicket(c echo.Context) error {
	// get id of wanted document
	id := c.Param("id")

	// call the usecase
	foundTicket, err := t.ticketUC.ShowTicket(c.Request().Context(), id)
	if err != nil {
		return c.JSON(rest.ErrorResponse(err))
	}

	return c.JSON(http.StatusOK, foundTicket)
}

func (t *ticketHandler) AvailableTickets(c echo.Context) error {
	// call the usecase
	tickets, err := t.ticketUC.AvailableTickets(c.Request().Context())
	if err != nil {
		return c.JSON(rest.ErrorResponse(err))
	}

	return c.JSON(http.StatusOK, tickets)
}
