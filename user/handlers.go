package user

import (
	"dona_tutti_api/middleware"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all user routes
func RegisterRoutes(g *echo.Group, service Service) {
	handler := NewHandler(service)

	// Auth routes
	authGroup := g.Group("/auth")
	authGroup.POST("/register", handler.Register)
	authGroup.POST("/login", handler.Login)
	authGroup.POST("/password-reset/request", handler.RequestPasswordReset)
	authGroup.POST("/password-reset/reset", handler.ResetPassword)

	// User routes
	userGroup := g.Group("/users")
	userGroup.GET("", handler.ListUsers)
	userGroup.GET("/:id", handler.GetUser)
	userGroup.POST("", handler.CreateUser)
	userGroup.PUT("/:id", handler.UpdateUser)
	userGroup.PUT("/:id/password", handler.UpdatePassword)

	// Protected user routes (require authentication)
	protectedUserGroup := userGroup.Group("", middleware.RequireAuth())
	protectedUserGroup.GET("/me", handler.GetMe)
}

// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterDTO true "User registration details"
// @Success 201 {object} RegisterResponseDTO
// @Failure 400 {object} errors.APIError
// @Router /auth/register [post]
func (h *Handler) Register(c echo.Context) error {
	var req RegisterDTO
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	id, err := h.service.Register(c.Request().Context(), req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, RegisterResponseDTO{ID: id})
}

// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginDTO true "Login credentials"
// @Success 200 {object} LoginResponseDTO
// @Failure 400 {object} errors.APIError
// @Router /auth/login [post]
func (h *Handler) Login(c echo.Context) error {
	var req LoginDTO
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	token, err := h.service.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, LoginResponseDTO{Token: *token})
}

// @Summary List all users
// @Description Get a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} GetUserResponseDTO
// @Failure 400 {object} errors.APIError
// @Router /users [get]
func (h *Handler) ListUsers(c echo.Context) error {
	users, err := h.service.ListUsers(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var response []GetUserResponseDTO
	for _, user := range users {
		response = append(response, GetUserResponseDTO{User: user})
	}
	return c.JSON(http.StatusOK, response)
}

// @Summary Get user by ID
// @Description Get user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} GetUserResponseDTO
// @Failure 400 {object} errors.APIError
// @Router /users/{id} [get]
func (h *Handler) GetUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	user, err := h.service.GetUser(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, GetUserResponseDTO{User: user})
}

// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body RegisterDTO true "User details"
// @Success 201 {object} RegisterResponseDTO
// @Failure 400 {object} errors.APIError
// @Router /users [post]
func (h *Handler) CreateUser(c echo.Context) error {
	var dto RegisterDTO
	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	id, err := h.service.CreateUser(c.Request().Context(), dto)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, RegisterResponseDTO{ID: id})
}

// @Summary Update user details
// @Description Update user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body UpdateUserDTO true "User details"
// @Success 200
// @Failure 400 {object} errors.APIError
// @Router /users/{id} [put]
func (h *Handler) UpdateUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	var dto UpdateUserDTO
	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	err = h.service.UpdateUser(c.Request().Context(), id, dto)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// @Summary Update user password
// @Description Update user's password
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param passwords body UpdatePasswordDTO true "Password update details"
// @Success 200
// @Failure 400 {object} errors.APIError
// @Router /users/{id}/password [put]
func (h *Handler) UpdatePassword(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	var req UpdatePasswordDTO
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	err = h.service.UpdatePassword(c.Request().Context(), id, req.CurrentPassword, req.NewPassword)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// @Summary Request password reset
// @Description Request a password reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param email body RequestPasswordResetDTO true "Email for password reset"
// @Success 200
// @Failure 400 {object} errors.APIError
// @Router /auth/password-reset/request [post]
func (h *Handler) RequestPasswordReset(c echo.Context) error {
	var req RequestPasswordResetDTO
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	err := h.service.RequestPasswordReset(c.Request().Context(), req.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// @Summary Reset password
// @Description Reset password using token
// @Tags auth
// @Accept json
// @Produce json
// @Param reset body ResetPasswordDTO true "Password reset details"
// @Success 200
// @Failure 400 {object} errors.APIError
// @Router /auth/password-reset/reset [post]
func (h *Handler) ResetPassword(c echo.Context) error {
	var req ResetPasswordDTO
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	err := h.service.ResetPassword(c.Request().Context(), req.Token, req.NewPassword)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// @Summary Get current user details
// @Description Get current authenticated user's details including role information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MeResponseDTO
// @Failure 401 {object} errors.APIError
// @Failure 500 {object} errors.APIError
// @Router /users/me [get]
func (h *Handler) GetMe(c echo.Context) error {
	// Get user ID from context (set by auth middleware)
	userIDValue := c.Get("user_id")
	if userIDValue == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	userID, err := uuid.Parse(userIDValue.(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID")
	}

	userMe, err := h.service.GetMe(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user information")
	}

	return c.JSON(http.StatusOK, userMe)
}
