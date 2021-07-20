package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/godocompany/livechat-api/models"
	"github.com/godocompany/livechat-api/utils"
	"gorm.io/gorm"
)

type AuthTokensService struct {
	DB            *gorm.DB
	SigningPepper string
}

// getSigningSecretKey gets the secret key used to sign and verify JWT tokens. The secret key combines
// the account's salt, the platform's salt, as well as the signing pepper for this server overall.
// Changes to any of those three values will result in possibly many tokens becoming invalid.
func (s *AuthTokensService) getSigningSecretKey(secret string) []byte {
	return []byte(utils.Sha256Hex(secret + s.SigningPepper))
}

// CreateToken creates an auth token for an account
func (s *AuthTokensService) CreateToken(
	account *models.Account,
	created time.Time,
	expire time.Time,
) (string, error) {

	// Create the claims for the token
	claims := jwt.MapClaims{
		"uid": account.ID,
		"cre": created.UTC().Unix(),
		"exp": expire.UTC().Unix(),
	}

	// Generate the token object
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

	// Sign the token to a string
	return token.SignedString(s.getSigningSecretKey(account.PasswordSalt))

}

// GetAccountForToken gets the account of the provided token string
func (s *AuthTokensService) GetAccountForToken(token string) (*models.Account, error) {

	// Decode the token
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		// Try to get the account of the token object
		account, err := s.getAccountFromTokenObj(token)
		if err != nil {
			fmt.Println("Err: ", err.Error())
			return nil, err
		}
		if account == nil {
			return nil, errors.New("no account matches token")
		}

		// Get the signing secret for the account
		return s.getSigningSecretKey(account.PasswordSalt), nil

	})
	if err != nil {
		return nil, err
	}

	// Get the object from the token
	return s.getAccountFromTokenObj(tokenObj)

}

func (s *AuthTokensService) validateTokenClaims(token *jwt.Token) (jwt.MapClaims, error) {

	// Cast the claims to a MapClaims instance
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("unrecognized claims type in token")
	}

	// Get the expiration date
	// It's a float64 because it's stores as a JSON number
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("token claims missing \"exp\" field")
	}
	if int64(exp) < time.Now().UTC().Unix() {
		return nil, errors.New("token has expired")
	}

	// Get the account UUID value
	if _, ok = claims["uid"]; !ok {
		return nil, errors.New("token claims missing \"uid\" field")
	}

	// Return no issues
	return claims, nil

}

// getAccountFromTokenObj gets the account from the provided token object
func (s *AuthTokensService) getAccountFromTokenObj(token *jwt.Token) (*models.Account, error) {

	// Validate the claims
	claims, err := s.validateTokenClaims(token)
	if err != nil {
		return nil, err
	}

	// Search the database for the model
	var account models.Account
	err = s.DB.
		Where("id = ?", claims["uid"]).
		First(&account).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil

}
