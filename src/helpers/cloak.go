package helpers

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

type CloakHelper struct {
	Client       *gocloak.GoCloak
	Realm        string
	ClientID     string
	ClientSecret string
}

func NewCloakHelper(url, realm, clientID, secret string) *CloakHelper {
	return &CloakHelper{
		Client:       gocloak.NewClient(url),
		Realm:        realm,
		ClientID:     clientID,
		ClientSecret: secret,
	}
}

// GetUserProfile fetches a user by ID using the Service Account token
func (s *CloakHelper) GetUserProfile(ctx context.Context, targetUserID string) (*gocloak.User, error) {
	// 1. Login as Service Account (Client Credentials)
	// Optimization Note: You should cache this token and only refresh when it expires.
	token, err := s.Client.LoginClient(ctx, s.ClientID, s.ClientSecret, s.Realm)
	if err != nil {
		return nil, fmt.Errorf("failed to login as service account: %w", err)
	}

	// 2. Use the Admin token to fetch the specific user
	user, err := s.Client.GetUserByID(ctx, token.AccessToken, s.Realm, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %s: %w", targetUserID, err)
	}

	return user, nil
}

func (s *CloakHelper) GetPublicUserData(ctx context.Context, targetUserID string) (map[string]interface{}, error) {
	user, err := s.GetUserProfile(ctx, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %s: %w", targetUserID, err)
	}

	// 3. Filter and return safe public data
	return map[string]interface{}{
		"id":        safeStr(user.ID),
		"username":  safeStr(user.Username),
		"firstName": safeStr(user.FirstName),
		"lastName":  safeStr(user.LastName),
		"email":     safeStr(user.Email),
	}, nil
}

// Helper to safely dereference string pointers
func safeStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
