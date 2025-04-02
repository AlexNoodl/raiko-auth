package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github/alexnoodl/raiko-auth/internal/models"
	"github/alexnoodl/raiko-auth/internal/service"
	"net/http"
)

type AuthHandler struct {
	authService *service.AuthService
	logger      *logrus.Logger
}

func NewAuthHandler(authService *service.AuthService, logger *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// Register
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя с указанными email, username и паролем
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} models.SuccessResponse "Пользователь успешно создан"
// @Failure 400 {object} models.ErrorResponse "Неверные данные или пароль не соответствует требованиям"
// @Failure 409 {object} models.ErrorResponse "Email или username уже занят"
// @Failure 500 {object} models.ErrorResponse "Ошибка сервера"
// @Router /api/v1/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	h.logger.Info("Received registration request")

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.WithError(err).Error("Failed to parse registration request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if user.Email == "" || user.Username == "" || user.Password == "" {
		h.logger.Error("Received empty required fields in registration request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Все поля (email, username, password) обязательны"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"email":    user.Email,
		"username": user.Username,
	}).Debug("Processing registration request")

	err := h.authService.Register(&user)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"email": user.Email,
			"error": err,
		}).Error("Registration failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithField("email", user.Email).Info("Registration request completed successfully")
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login
// @Summary Аутентификация пользователя
// @Description Аутентифицирует пользователя и возвращает JWT-токен
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.LoginRequest true "Данные для входа"
// @Success 200 {object} models.LoginResponse "Успешный вход"
// @Failure 400 {object} models.ErrorResponse "Неверные данные"
// @Failure 500 {object} models.ErrorResponse "Ошибка сервера"
// @Router /api/v1/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	h.logger.Info("Received login request")

	var credentials struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		h.logger.WithError(err).Error("Failed to parse login request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	h.logger.WithField("identifier", credentials.Login).Debug("Processing login request")

	token, err := h.authService.Login(credentials.Login, credentials.Password)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"identifier": credentials.Login,
			"error":      err,
		}).Error("Login failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithField("login", credentials.Login).Info("Login request completed successfully")
	c.JSON(http.StatusOK, gin.H{"token": token})
}
