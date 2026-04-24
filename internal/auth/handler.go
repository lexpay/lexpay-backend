package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	pgtype "github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/lexpay/internal/db"
	"github.com/luponetn/lexpay/internal/utils"
)

type Handler struct {
	service Svc
}

func NewHandler(service Svc) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleSignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid request payload",
			"error":   err.Error(),
		})
		slog.Error("invalid request payload", "error", err.Error())
		return
	}

	args := db.CreateUserOnSignupParams{
		Name:         req.Name,
		Email:        req.Email,
		PhoneNumber:  pgtype.Text{String: req.PhoneNumber, Valid: true},
		Nationality:  req.Nationality,
		PasswordHash: req.Password,
	}

	user, err := h.service.SignUp(c.Request.Context(), args)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "failed",
				"message": "a user with this email already exists",
			})
			slog.Warn("signup failed: user already exists", "email", req.Email)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "failed to create user",
			"error":   err.Error(),
		})
		slog.Error("failed to sign up user", "error", err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "user created successfully",
		"data":    user,
	})
}

func (h *Handler) HandleSignIn(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid request payload",
			"error":   err.Error(),
		})
		return
	}

	tokens, err := h.service.SignIn(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "invalid email or password",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "failed to sign in",
			"error":   err.Error(),
		})
		return
	}

	// Set access token cookie — short-lived (20 minutes)
	// HttpOnly: true prevents JavaScript from reading it (XSS protection)
	// Secure: set to true in production so it is only sent over HTTPS
	c.SetCookie("access_token", tokens.AccessToken, int(utils.AccessTokenDuration.Seconds()), "/", "", false, true)

	// Set refresh token cookie — long-lived (7 days)
	// Path is restricted to /auth/refresh so it is only sent when refreshing
	c.SetCookie("refresh_token", tokens.RefreshToken, int(utils.RefreshTokenDuration.Seconds()), "/auth/refresh", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "signed in successfully",
	})
}

// HandleRefreshToken reads the refresh token from the cookie, validates it,
// and issues a brand-new access + refresh token pair.
func (h *Handler) HandleRefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"message": "refresh token not found",
		})
		return
	}

	tokens, err := h.service.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "invalid or expired refresh token",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "failed to refresh token",
			"error":   err.Error(),
		})
		return
	}

	// Set the new access token cookie
	c.SetCookie("access_token", tokens.AccessToken, int(utils.AccessTokenDuration.Seconds()), "/", "", false, true)

	// Set the new refresh token cookie (token rotation)
	c.SetCookie("refresh_token", tokens.RefreshToken, int(utils.RefreshTokenDuration.Seconds()), "/auth/refresh", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "tokens refreshed successfully",
	})
}