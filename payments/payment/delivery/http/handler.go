package http

import (
	"github.com/halilylm/gommon/middlewares"
	"github.com/halilylm/gommon/rest"
	"github.com/halilylm/secondhand/payments/payment/usecase"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type paymentHandler struct {
	paymentUC usecase.Payment
}

// NewPaymentHandler handler for payment
func NewPaymentHandler(g *echo.Group, paymentUC usecase.Payment) {
	handler := &paymentHandler{paymentUC: paymentUC}

	// jwt middleware
	g.Use(middlewares.CurrentUser("jwt"))

	g.POST("/:order_id", handler.NewPayment)
}

func (p *paymentHandler) NewPayment(c echo.Context) error {
	// id of wanted order
	orderID := c.Param("order_id")

	// get user from the context
	user := middlewares.UserFromContext(c)
	log.Println(user.ID)

	// call the usecase
	createdOrder, err := p.paymentUC.Pay(c.Request().Context(), orderID)
	if err != nil {
		return c.JSON(rest.ErrorResponse(err))
	}

	return c.JSON(http.StatusCreated, createdOrder)
}
