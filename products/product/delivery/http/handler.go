package http

import (
	"github.com/halilylm/gommon/middlewares"
	"github.com/halilylm/gommon/rest"
	"github.com/halilylm/gommon/utils"
	"github.com/halilylm/secondhand/product/domain"
	"github.com/halilylm/secondhand/product/product/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
)

type productHandler struct {
	productUC usecase.Product
}

// NewProductHandler handler for products
func NewProductHandler(g *echo.Group, productUC usecase.Product) {
	handler := &productHandler{productUC: productUC}

	// jwt middleware
	g.Use(middlewares.CurrentUser("jwt"))

	g.POST("/", handler.NewProduct)
	g.PUT("/:id", handler.UpdateProduct)
	g.GET("/:id", handler.ShowProduct)
	g.GET("/", handler.AvailableProducts)
}

func (p *productHandler) NewProduct(c echo.Context) error {
	var product domain.Product

	// bind request body to product
	if err := c.Bind(&product); err != nil {
		return c.JSON(rest.ErrorResponse(rest.NewBadRequestError(err.Error())))
	}

	// validate the struct
	if err := utils.ValidateStruct(&product); err != nil {
		return c.JSON(rest.ErrorResponse(rest.NewValidationErrors(err)))
	}
	
	// get user from the context
	user := middlewares.UserFromContext(c)
	product.UserID = user.ID

	// call the usecase
	createdProduct, err := p.productUC.NewProduct(c.Request().Context(), &product)
	if err != nil {
		// errors returning from usecase layer will be rest errors
		// so err can be used directly
		return c.JSON(rest.ErrorResponse(err))
	}

	return c.JSON(http.StatusCreated, createdProduct)
}

func (p *productHandler) UpdateProduct(c echo.Context) error {
	var product domain.Product

	// bind request body to ticket
	if err := c.Bind(&product); err != nil {
		return c.JSON(rest.ErrorResponse(rest.NewBadRequestError(err.Error())))
	}

	// validate the struct
	if err := utils.ValidateStruct(&product); err != nil {
		return c.JSON(rest.ErrorResponse(rest.NewValidationErrors(err)))
	}

	// fill the ticket id
	product.ID = c.Param("id")

	// call the usecase
	updatedProduct, err := p.productUC.UpdateProduct(c.Request().Context(), &product)
	if err != nil {
		// errors returning from usecase layer will be rest errors
		// so err can be used directly
		return c.JSON(rest.ErrorResponse(err))
	}
	return c.JSON(http.StatusOK, updatedProduct)
}

func (p *productHandler) ShowProduct(c echo.Context) error {
	// id of wanted product
	id := c.Param("id")

	// call the usecase
	foundProduct, err := p.productUC.ShowProduct(c.Request().Context(), id)
	if err != nil {
		return c.JSON(rest.ErrorResponse(err))
	}

	return c.JSON(http.StatusOK, foundProduct)
}

func (p *productHandler) AvailableProducts(c echo.Context) error {
	// call the usecase
	products, err := p.productUC.AvailableProducts(c.Request().Context())
	if err != nil {
		return c.JSON(rest.ErrorResponse(err))
	}

	return c.JSON(http.StatusOK, products)
}
