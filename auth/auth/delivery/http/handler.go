package http

import (
	"net/http"

	"github.com/halilylm/gommon/middlewares"
	"github.com/halilylm/gommon/rest"
	"github.com/halilylm/gommon/utils"
	"github.com/halilylm/secondhand/auth/auth/usecase"
	"github.com/halilylm/secondhand/auth/domain"
	"github.com/labstack/echo/v4"
)

type authHandler struct {
	authUC usecase.Auth
}

// NewAuthHandler handler for auth
func NewAuthHandler(g *echo.Group, authUC usecase.Auth) {
	handler := &authHandler{authUC: authUC}
	g.POST("/signup", handler.SignUp)
	g.GET("/signout", handler.SignOut)
	g.POST("/signin", handler.SignIn)

	// jwt middleware
	g.Use(middlewares.CurrentUser("jwt"))

	g.GET("/currentuser", handler.CurrentUser)
}

// SignUp the user
func (a *authHandler) SignUp(c echo.Context) error {
	var user domain.User

	// bind request body to user
	if err := c.Bind(&user); err != nil {
		return c.JSON(rest.ErrorResponse(rest.NewBadRequestError(err.Error())))
	}

	// validate the struct
	if err := utils.ValidateStruct(&user); err != nil {
		return c.JSON(rest.ErrorResponse(rest.NewValidationErrors(err)))
	}

	// call the usecase
	createdUser, err := a.authUC.SignUp(c.Request().Context(), &user)
	if err != nil {
		// errors returning from usecase layer will be rest errors
		// so err can be used directly
		return c.JSON(rest.ErrorResponse(err))
	}

	// generate jwt token
	token, err := utils.GenerateJWTToken(createdUser.Email, createdUser.ID)
	if err != nil {
		return c.JSON(rest.ErrorResponse(rest.NewInternalServerError()))
	}

	// set the cookie
	setCookie(c, token)

	return c.JSON(http.StatusCreated, createdUser)
}

// SignIn the user
func (a *authHandler) SignIn(c echo.Context) error {
	var user domain.User

	// bind request body to user
	if err := c.Bind(&user); err != nil {
		return c.JSON(rest.ErrorResponse(rest.NewBadRequestError(err.Error())))
	}

	// validate the struct
	if err := utils.ValidateStruct(&user); err != nil {
		return c.JSON(rest.ErrorResponse(rest.NewValidationErrors(err)))
	}

	// call the usecase
	loginUser, err := a.authUC.SignIn(c.Request().Context(), &user)
	if err != nil {
		// errors returning from usecase layer will be rest errors
		// so err can be used directly
		return c.JSON(rest.ErrorResponse(err))
	}

	// generate jwt token
	token, err := utils.GenerateJWTToken(loginUser.Email, loginUser.ID)
	if err != nil {
		return c.JSON(rest.ErrorResponse(rest.NewInternalServerError()))
	}

	// set the cookie
	setCookie(c, token)

	return c.JSON(http.StatusCreated, loginUser)
}

// CurrentUser returns the current user
func (a *authHandler) CurrentUser(c echo.Context) error {
	return c.JSON(http.StatusOK, middlewares.UserFromContext(c))
}

// SignOut the user
func (a *authHandler) SignOut(c echo.Context) error {
	unsetCookie(c)
	return c.NoContent(http.StatusOK)
}

// setCookie sets jwt token in cookie
func setCookie(c echo.Context, token string) {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
	}
	c.SetCookie(cookie)
}

// unsetCookie unsets the cookie
func unsetCookie(c echo.Context) {
	cookie := &http.Cookie{
		Name:   "jwt",
		MaxAge: -1,
	}
	c.SetCookie(cookie)
}
