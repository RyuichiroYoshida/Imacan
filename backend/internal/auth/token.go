package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrInvalidToken           = errors.New("invalid token")
	ErrExpiredToken           = errors.New("expired token")
	ErrDiscordNotConfigured   = errors.New("discord oauth is not configured")
	ErrDiscordAuthentication  = errors.New("discord authentication failed")
	ErrDiscordUserUnavailable = errors.New("discord user unavailable")
)

type Claims struct {
	UserID    string `json:"user_id"`
	ExpiresAt int64  `json:"exp"`
}

type Service struct {
	secret     []byte
	ttl        time.Duration
	discord    DiscordConfig
	httpClient *http.Client
}

type DiscordConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	TokenURL     string
	UserURL      string
}

func NewService(secret string, ttl time.Duration) *Service {
	return &Service{
		secret:     []byte(secret),
		ttl:        ttl,
		httpClient: http.DefaultClient,
	}
}

func (s *Service) ConfigureDiscord(config DiscordConfig, httpClient *http.Client) {
	if config.TokenURL == "" {
		config.TokenURL = "https://discord.com/api/oauth2/token"
	}
	if config.UserURL == "" {
		config.UserURL = "https://discord.com/api/users/@me"
	}
	s.discord = config
	if httpClient != nil {
		s.httpClient = httpClient
	}
}

func (s *Service) AuthenticateDiscord(ctx context.Context, code string, redirectURI *string) (string, error) {
	if code == "" {
		return "", ErrDiscordAuthentication
	}
	if s.discord.ClientID == "" || s.discord.ClientSecret == "" {
		return "", ErrDiscordNotConfigured
	}

	token, err := s.exchangeDiscordCode(ctx, code, redirectURI)
	if err != nil {
		return "", err
	}

	discordID, err := s.fetchDiscordUserID(ctx, token)
	if err != nil {
		return "", err
	}
	return "discord:" + discordID, nil
}

func (s *Service) Issue(userID string, now time.Time) (string, int32, error) {
	if userID == "" {
		return "", 0, errors.New("user id is required")
	}

	expiresAt := now.Add(s.ttl)
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	claims := Claims{
		UserID:    userID,
		ExpiresAt: expiresAt.Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", 0, err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", 0, err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJSON)
	signingInput := encodedHeader + "." + encodedClaims
	signature := s.sign(signingInput)

	return signingInput + "." + signature, int32(s.ttl.Seconds()), nil
}

func (s *Service) Verify(token string, now time.Time) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, ErrInvalidToken
	}

	signingInput := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(s.sign(signingInput)), []byte(parts[2])) {
		return Claims{}, ErrInvalidToken
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	var claims Claims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return Claims{}, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	if claims.UserID == "" {
		return Claims{}, ErrInvalidToken
	}
	if now.Unix() >= claims.ExpiresAt {
		return Claims{}, ErrExpiredToken
	}

	return claims, nil
}

func (s *Service) sign(input string) string {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(input))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (s *Service) exchangeDiscordCode(ctx context.Context, code string, redirectURI *string) (string, error) {
	form := url.Values{}
	form.Set("client_id", s.discord.ClientID)
	form.Set("client_secret", s.discord.ClientSecret)
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	if redirectURI != nil && *redirectURI != "" {
		form.Set("redirect_uri", *redirectURI)
	} else if s.discord.RedirectURI != "" {
		form.Set("redirect_uri", s.discord.RedirectURI)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.discord.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")

	response, err := s.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrDiscordAuthentication, err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("%w: token endpoint returned %d", ErrDiscordAuthentication, response.StatusCode)
	}

	var body struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("%w: %v", ErrDiscordAuthentication, err)
	}
	if body.AccessToken == "" {
		return "", ErrDiscordAuthentication
	}
	return body.AccessToken, nil
}

func (s *Service) fetchDiscordUserID(ctx context.Context, accessToken string) (string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, s.discord.UserURL, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Accept", "application/json")

	response, err := s.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrDiscordUserUnavailable, err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		_, _ = io.Copy(io.Discard, response.Body)
		return "", fmt.Errorf("%w: user endpoint returned %d", ErrDiscordUserUnavailable, response.StatusCode)
	}

	var body struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("%w: %v", ErrDiscordUserUnavailable, err)
	}
	if body.ID == "" {
		return "", ErrDiscordUserUnavailable
	}
	return body.ID, nil
}
