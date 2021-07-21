package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/godocompany/livechat-api/models"
	"gorm.io/gorm"
)

type ChatUserInfo struct {
	Username  string `json:"username"`
	IpAddress string `json:"ip_address"`
}

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
	user *ChatUserInfo,
	untilDate *time.Time,
) (*models.MutedUser, error) {

	// If the user info is missing both fields
	if len(user.Username) == 0 && len(user.IpAddress) == 0 {
		return nil, nil
	}

	// Create the until date
	var until sql.NullTime
	if untilDate != nil {
		until = sql.NullTime{
			Valid: true,
			Time:  *untilDate,
		}
	}

	// Add an entry to mute the user
	mutedUser := models.MutedUser{
		OrganizationID: organizationID,
		UntilDate:      until,
		CreatedDate:    time.Now(),
	}
	if len(user.Username) > 0 {
		mutedUser.Username = sql.NullString{
			Valid:  true,
			String: user.Username,
		}
	}
	if len(user.IpAddress) > 0 {
		mutedUser.IpAddress = sql.NullString{
			Valid:  true,
			String: user.IpAddress,
		}
	}
	if err := s.DB.Create(&mutedUser).Error; err != nil {
		return nil, err
	}
	return &mutedUser, nil

}

func (s *ChatService) UnmuteUser(
	organizationID uint64,
	user *ChatUserInfo,
) error {

	// If the user info is missing both fields
	if len(user.Username) == 0 && len(user.IpAddress) == 0 {
		return nil
	}

	// Construct the query
	query := s.DB.
		Model(&models.MutedUser{}).
		Where("deleted_date IS NULL").
		Where("organization_id = ?", organizationID)

	// Create the ors
	ors := s.DB

	// Add the username and/or IP address
	if len(user.Username) > 0 {
		ors = ors.Or("username LIKE ?", user.Username)
	}
	if len(user.IpAddress) > 0 {
		ors = ors.Or("ip_address LIKE ?", user.IpAddress)
	}

	// Update all of the muted users and mark as deleted
	return query.
		Where(ors).
		Update("deleted_date", time.Now()).
		Error

}

func (s *ChatService) IsUserMuted(
	organizationID uint64,
	user *ChatUserInfo,
) (bool, error) {

	// If the user info is missing both fields
	if len(user.Username) == 0 && len(user.IpAddress) == 0 {
		return false, nil
	}

	// Construct the query
	query := s.DB.
		Where("deleted_date IS NULL").
		Where("until_date IS NULL OR until_date < ?", time.Now()).
		Where("organization_id = ?", organizationID)

	// Add the username and/or IP address
	if len(user.Username) > 0 {
		query = query.Where("username LIKE ?", user.Username)
	}
	if len(user.IpAddress) > 0 {
		query = query.Where("ip_address LIKE ?", user.IpAddress)
	}

	var mutedUser models.MutedUser
	if err := query.First(&mutedUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetBannedWords gets all of the banned words for an organization. The slice returned also includes all of the
// platform-wide banned words
func (s *ChatService) GetBannedWords(organizationID uint64) ([]*models.BannedWord, error) {
	var bannedWords []*models.BannedWord
	err := s.DB.
		Where("deleted_date IS NULL").
		Where("organization_id IS NULL OR organization_id = ?", organizationID).
		Find(&bannedWords).
		Error
	if err != nil {
		return nil, err
	}
	return bannedWords, nil
}

func (s *ChatService) checkMessageAgainstBannedWord(message string, bw *models.BannedWord) bool {

	return true
}

// CanSendMessage determines if a given message can be sent from a user to a chatroom
func (s *ChatService) CanSendMessage(
	chatRoom *models.ChatRoom,
	user *ChatUserInfo,
	message string,
) (bool, *models.BannedWord, error) {

	// Check if the user is banned
	muted, err := s.IsUserMuted(chatRoom.OrganizationID, user)
	if err != nil {
		return false, nil, err
	}
	if muted {
		return false, nil, nil
	}

	// Check for all the banned words
	bannedWords, err := s.GetBannedWords(chatRoom.OrganizationID)
	if err != nil {
		return false, nil, err
	}

	// Loop through the banned words
	for _, bw := range bannedWords {
		if !s.checkMessageAgainstBannedWord(message, bw) {
			return false, bw, nil
		}
	}

	// The message looks good
	return true, nil, nil

}
