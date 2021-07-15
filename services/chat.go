package services

import (
	"errors"
	"time"

	"github.com/godocompany/livechat-api/models"
	"gorm.io/gorm"
)

// ChatService manages chat moderation
type ChatService struct {
	DB *gorm.DB
}

// GetChatRoomByIdentifier gets the chat room with the provided identifier
func (s *ChatService) GetChatRoomByIdentifier(identifier string) (*models.ChatRoom, error) {
	var chatRoom models.ChatRoom
	err := s.DB.
		Where("deleted_date IS NULL").
		Where("identifier = ?", identifier).
		First(&chatRoom).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &chatRoom, nil
}

func (s *ChatService) MuteUser(
	organizationID uint64,
	username string,
) (*models.MutedUser, error) {

	// Add an entry to mute the user
	mutedUser := models.MutedUser{
		Username:       username,
		OrganizationID: organizationID,
		CreatedDate:    time.Now(),
	}
	if err := s.DB.Create(&mutedUser).Error; err != nil {
		return nil, err
	}
	return &mutedUser, nil

}

func (s *ChatService) UnmuteUser(
	organizationID uint64,
	username string,
) error {
	return s.DB.
		Model(&models.MutedUser{}).
		Where("deleted_date IS NULL").
		Where("organization_id = ?", organizationID).
		Where("username LIKE ?", username).
		Update("deleted_date", time.Now()).
		Error
}

func (s *ChatService) IsUserMuted(
	organizationID uint64,
	username string,
) (bool, error) {
	var mutedUser models.MutedUser
	err := s.DB.
		Where("deleted_date IS NULL").
		Where("username LIKE ?", username).
		Where("organization_id = ?", organizationID).
		First(&mutedUser).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
