package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"log"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type NotificationService struct {
	repo        ports.NotificationRepository
	client      *messaging.Client
	firebaseApp *firebase.App
}

func NewNotificationService(repo ports.NotificationRepository, serviceAccountPath string) (*NotificationService, error) {
	opt := option.WithCredentialsFile(serviceAccountPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, err
	}

	return &NotificationService{
		repo:        repo,
		client:      client,
		firebaseApp: app,
	}, nil
}

func (s *NotificationService) SendNotification(ctx context.Context, userID uuid.UUID, title, body string, notifType domain.NotificationType, data map[string]string) error {
	// Get user's FCM tokens
	tokens, err := s.repo.GetUserTokens(ctx, userID)
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		log.Printf("No FCM tokens found for user %s", userID)
		return nil // Not an error, user just doesn't have tokens registered
	}

	// Prepare notification data
	if data == nil {
		data = make(map[string]string)
	}
	data["type"] = string(notifType)
	data["user_id"] = userID.String()

	// Send to all user's tokens
	tokenStrings := make([]string, len(tokens))
	for i, token := range tokens {
		tokenStrings[i] = token.Token
	}

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data:   data,
		Tokens: tokenStrings,
	}

	response, err := s.client.SendMulticast(ctx, message)
	if err != nil {
		return err
	}

	// Handle failed tokens (expired or invalid)
	if response.FailureCount > 0 {
		for idx, resp := range response.Responses {
			if !resp.Success {
				// Token is invalid, remove it
				log.Printf("Failed to send to token %s: %v", tokenStrings[idx], resp.Error)
				_ = s.repo.DeleteToken(ctx, tokenStrings[idx])
			}
		}
	}

	// Save notification to history
	notification := &domain.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     title,
		Message:   body,
		Type:      string(notifType),
		Link:      "", // TODO: Add link support
		Read:      false,
		CreatedAt: time.Now(),
	}

	return s.repo.Save(ctx, notification)
}

func (s *NotificationService) SendBatchNotification(ctx context.Context, userIDs []uuid.UUID, title, body string, notifType domain.NotificationType, data map[string]string) error {
	for _, userID := range userIDs {
		// Send asynchronously
		go func(uid uuid.UUID) {
			if err := s.SendNotification(context.Background(), uid, title, body, notifType, data); err != nil {
				log.Printf("Failed to send notification to user %s: %v", uid, err)
			}
		}(userID)
	}
	return nil
}

func (s *NotificationService) RegisterFCMToken(ctx context.Context, userID uuid.UUID, token, deviceType string) error {
	fcmToken := &domain.FCMToken{
		ID:         uuid.New(),
		UserID:     userID,
		Token:      token,
		DeviceType: deviceType,
		CreatedAt:  time.Now(),
		LastUsedAt: time.Now(),
	}

	return s.repo.SaveToken(ctx, fcmToken)
}

func (s *NotificationService) UnregisterFCMToken(ctx context.Context, userID uuid.UUID, token string) error {
	return s.repo.DeleteToken(ctx, token)
}

func (s *NotificationService) GetHistory(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error) {
	return s.repo.FindByUserID(ctx, userID, limit)
}

// NoOpNotificationService implementation for when Firebase is not configured
type NoOpNotificationService struct{}

func NewNoOpNotificationService() *NoOpNotificationService {
	return &NoOpNotificationService{}
}

func (s *NoOpNotificationService) SendNotification(ctx context.Context, userID uuid.UUID, title, body string, notifType domain.NotificationType, data map[string]string) error {
	log.Printf("[NoOp] SendNotification to %s: %s - %s", userID, title, body)
	return nil
}

func (s *NoOpNotificationService) SendBatchNotification(ctx context.Context, userIDs []uuid.UUID, title, body string, notifType domain.NotificationType, data map[string]string) error {
	log.Printf("[NoOp] SendBatchNotification to %d users: %s - %s", len(userIDs), title, body)
	return nil
}

func (s *NoOpNotificationService) RegisterFCMToken(ctx context.Context, userID uuid.UUID, token, deviceType string) error {
	log.Printf("[NoOp] RegisterFCMToken for %s: %s", userID, token)
	return nil
}

func (s *NoOpNotificationService) UnregisterFCMToken(ctx context.Context, userID uuid.UUID, token string) error {
	log.Printf("[NoOp] UnregisterFCMToken for %s: %s", userID, token)
	return nil
}

func (s *NoOpNotificationService) GetHistory(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error) {
	return []domain.Notification{}, nil
}
