package services

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/suipic/backend/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	dbService  *DatabaseService
	jwtSecret  string
	jwtExpiry  string
	adminEmail string
	adminPass  string
	adminUser  string
}

type JWTClaims struct {
	UserID   int64           `json:"user_id"`
	Email    string          `json:"email"`
	Username string          `json:"username"`
	Role     models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(dbService *DatabaseService, jwtSecret, jwtExpiry, adminEmail, adminPass, adminUser string) (*AuthService, error) {
	service := &AuthService{
		dbService:  dbService,
		jwtSecret:  jwtSecret,
		jwtExpiry:  jwtExpiry,
		adminEmail: adminEmail,
		adminPass:  adminPass,
		adminUser:  adminUser,
	}

	if err := service.seedAdminUser(); err != nil {
		return nil, fmt.Errorf("failed to seed admin user: %w", err)
	}

	return service, nil
}

func (s *AuthService) seedAdminUser() error {
	if s.adminEmail == "" || s.adminPass == "" || s.adminUser == "" {
		return nil
	}

	existingUser, err := s.dbService.GetUserByEmail(s.adminEmail)
	if err != nil {
		return err
	}

	hashedPassword, err := s.HashPassword(s.adminPass)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	if existingUser != nil {
		// Update existing admin user's password if it doesn't match
		if err := s.CheckPassword(existingUser.PasswordHash, s.adminPass); err != nil {
			existingUser.PasswordHash = hashedPassword
			// Ensure role is admin
			existingUser.Role = models.RoleAdmin
			if err := s.dbService.GetUserRepo().Update(context.Background(), existingUser); err != nil {
				return fmt.Errorf("failed to update admin user: %w", err)
			}
			fmt.Println("Admin user password updated from environment configuration")
		}
		return nil
	}

	_, err = s.dbService.CreateUser(s.adminEmail, s.adminUser, hashedPassword, models.RoleAdmin)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	return nil
}

func (s *AuthService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (s *AuthService) CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	expiry, err := time.ParseDuration(s.jwtExpiry)
	if err != nil {
		expiry = 24 * time.Hour
	}

	claims := JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *AuthService) Register(email, username, password string, role models.UserRole) (*models.User, error) {
	return s.RegisterWithFriendlyName(email, username, password, "", role)
}

func (s *AuthService) RegisterWithFriendlyName(email, username, password, friendlyName string, role models.UserRole) (*models.User, error) {
	existingUser, err := s.dbService.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with email already exists")
	}

	existingUser, err = s.dbService.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with username already exists")
	}

	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := s.dbService.CreateUserWithFriendlyName(email, username, hashedPassword, friendlyName, role)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	user, err := s.dbService.GetUserByEmail(email)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	if err := s.CheckPassword(user.PasswordHash, password); err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) LoginWithUsernameOrEmail(username, email, password string) (*models.User, string, error) {
	var user *models.User
	var err error

	if username != "" {
		user, err = s.dbService.GetUserByUsername(username)
		if err != nil {
			return nil, "", err
		}
	}

	if user == nil && email != "" {
		user, err = s.dbService.GetUserByEmail(email)
		if err != nil {
			return nil, "", err
		}
	}

	if user == nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	if err := s.CheckPassword(user.PasswordHash, password); err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) GetUserByID(userID int64) (*models.User, error) {
	return s.dbService.GetUserByID(userID)
}

func (s *AuthService) GetUserByUsername(username string) (*models.User, error) {
	return s.dbService.GetUserByUsername(username)
}

func (s *AuthService) CreatePhotographerClient(photographerID, clientID int64) (*models.PhotographerClient, error) {
	return s.dbService.CreatePhotographerClient(photographerID, clientID)
}

func (s *AuthService) GetPhotographerClient(photographerID, clientID int64) (*models.PhotographerClient, error) {
	return s.dbService.GetPhotographerClient(photographerID, clientID)
}

func (s *AuthService) GetClientsByPhotographer(photographerID int64) ([]*models.User, error) {
	return s.dbService.GetClientsByPhotographer(photographerID)
}

func (s *AuthService) SearchClientsByUsername(username string) ([]*models.User, error) {
	return s.dbService.SearchClientsByUsername(username)
}

func (s *AuthService) GetClientPhotographerCounts(clients []*models.User) (map[int64]int, error) {
	return s.dbService.GetClientPhotographerCounts(clients)
}
