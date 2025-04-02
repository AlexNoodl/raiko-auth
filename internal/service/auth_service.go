package service

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github/alexnoodl/raiko-auth/internal/models"
	"github/alexnoodl/raiko-auth/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	db     *mongo.Database
	logger *logrus.Logger
	jwtKey []byte
}

func NewAuthService(db *mongo.Database, logger *logrus.Logger, jwtKey []byte) *AuthService {
	return &AuthService{
		db:     db,
		logger: logger,
		jwtKey: jwtKey,
	}
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func (s *AuthService) Register(user *models.User) error {
	s.logger.WithFields(logrus.Fields{
		"email":    user.Email,
		"username": user.Username,
	}).Info("Starting user registration")

	if !utils.IsValidPassword(user.Password) {
		s.logger.WithField("email", user.Email).Warn("Password validation failed")
		return errors.New("incorrect password type")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := s.db.Collection("users").CountDocuments(ctx, bson.M{"or": []bson.M{
		{"email": user.Email},
		{"username": user.Username},
	}})

	if err != nil {
		s.logger.WithError(err).Error("Failed to check existing users")
		return err
	}
	if count > 0 {
		s.logger.WithFields(logrus.Fields{
			"email":    user.Email,
			"username": user.Username,
		}).Warn("Email or username already exists")
		return errors.New("email or username already exists")
	}

	s.logger.WithField("email", user.Email).Debug("Hashing password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.WithError(err).Error("Failed to hash password")
		return err
	}
	// Todo нужно будет в дальнейшем isActive делать false до подтверждения регистрации по email
	user.Password = string(hashedPassword)
	user.IsActive = true

	s.logger.WithField("email", user.Email).Debug("Inserting user into database")
	_, err = s.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		s.logger.WithError(err).Error("Failed to insert user")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"email":    user.Email,
		"username": user.Username,
	}).Info("User registered successfully")
	return nil
}

func (s *AuthService) Login(login, password string) (string, error) {
	s.logger.WithField("login", login).Info("Starting login attempt")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User

	err := s.db.Collection("users").FindOne(ctx, bson.M{"$or": []bson.M{
		{"email": login},
		{"username": login},
	}}).Decode(&user)

	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"login": login,
			"error": err,
		}).Warn("User not found")
		return "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		s.logger.WithField("email", user.Email).Warn("Login attempt on inactive account")
		return "", errors.New("account is not active")
	}

	s.logger.WithField("email", user.Email).Debug("Verifying password")
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"email": user.Email,
			"error": err,
		}).Warn("Password verification failed")
		return "", errors.New("invalid credentials")
	}

	s.logger.WithField("email", user.Email).Debug("Generating JWT token")
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		s.logger.WithError(err).Error("Failed to generate JWT token")
		return "", err
	}

	s.logger.WithField("email", user.Email).Info("Login successful")
	return tokenString, nil
}
