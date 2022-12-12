package http

import (
	"github.com/halilylm/gommon/middlewares"
	"github.com/halilylm/gommon/rest"
	"github.com/halilylm/ticketing/orders/orders/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
)

type orderHandler struct {
	orderUC usecase.Order
}

// NewOrderHandler handler for order
func NewOrderHandler(g *echo.Group, orderUC usecase.Order) {
	handler := &orderHandler{orderUC: orderUC}

	// jwt middleware
	g.Use(middlewares.CurrentUser("jwt"))

	g.POST("/:ticket_id", handler.NewOrder)
	g.GET("/:id", handler.ShowOrder)
	g.DELETE("/:id", handler.DeleteOrder)
	g.GET("/", handler.ListOrders)
}

func (o *orderHandler) NewOrder(c echo.Context) error {
	// id of wanted order
	ticketID := c.Param("ticket_id")

	// get user from the context
	user := middlewares.UserFromContext(c)

	// call the usecase
	createdOrder, err := o.orderUC.NewOrder(c.Request().Context(), ticketID, user.ID)
	if err != nil {
		return c.JSON(rest.ErrorResponse(err))
	}

	return c.JSON(http.StatusCreated, createdOrder)
}

func (o *orderHandler) ShowOrder(c echo.Context) error {
	// id of wanted order
	id := c.Param("id")

	// get user from the context
	user := middlewares.UserFromContext(c)

	// call the usecase
	foundOrder, err := o.orderUC.ShowOrder(c.Request().Context(), id, user.ID)
	if err != nil {
		return c.JSON(rest.ErrorResponse(err))
	}

	return c.JSON(http.StatusOK, foundOrder)
}

func (o *orderHandler) DeleteOrder(c echo.Context) error {
	// id of wanted order
	id := c.Param("id")

	// get user from the context
	user := middlewares.UserFromContext(c)

	// call the usecase
	err := o.orderUC.DeleteOrder(c.Request().Context(), id, user.ID)
	if err != nil {
		return c.JSON(rest.ErrorResponse(err))
	}

	return c.NoContent(http.StatusNoContent)
}

func (o *orderHandler) ListOrders(c echo.Context) error {
	// get user from the context
	user := middlewares.UserFromContext(c)

	// call the usecase
	orders, err := o.orderUC.ListUserOrders(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(rest.ErrorResponse(err))
	}

	return c.JSON(http.StatusOK, orders)
}
