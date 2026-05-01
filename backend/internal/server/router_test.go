package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RyuichiroYoshida/imacan/backend/internal/auth"
	"github.com/RyuichiroYoshida/imacan/backend/internal/generated"
	"github.com/RyuichiroYoshida/imacan/backend/internal/presence"
)

func TestAuthPresenceFlow(t *testing.T) {
	authService := auth.NewService("test-secret", time.Hour)
	presenceService := presence.NewService(presence.NewMemoryStore(), 105*time.Minute, 20, 30)
	handler := NewHandler(authService, presenceService)
	router := NewRouter(handler, authService)

	token := requestToken(t, router)

	updateBody := []byte(`{"activity":"SELF_STUDY","lat":35.68,"lng":139.76}`)
	updateRequest := httptest.NewRequest(http.MethodPost, "/presence", bytes.NewReader(updateBody))
	updateRequest.Header.Set("Authorization", "Bearer "+token)
	updateRequest.Header.Set("Content-Type", "application/json")
	updateResponse := httptest.NewRecorder()
	router.ServeHTTP(updateResponse, updateRequest)

	if updateResponse.Code != http.StatusOK {
		t.Fatalf("expected update status 200, got %d: %s", updateResponse.Code, updateResponse.Body.String())
	}

	summaryRequest := httptest.NewRequest(http.MethodGet, "/presence/summary", nil)
	summaryResponse := httptest.NewRecorder()
	router.ServeHTTP(summaryResponse, summaryRequest)

	if summaryResponse.Code != http.StatusOK {
		t.Fatalf("expected summary status 200, got %d: %s", summaryResponse.Code, summaryResponse.Body.String())
	}

	var summary generated.PresenceSummaryResponseBody
	if err := json.NewDecoder(summaryResponse.Body).Decode(&summary); err != nil {
		t.Fatal(err)
	}
	if summary.Total != 1 || summary.SelfStudy != 1 || summary.Class != 0 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
}

func TestPresenceRequiresBearerToken(t *testing.T) {
	authService := auth.NewService("test-secret", time.Hour)
	presenceService := presence.NewService(presence.NewMemoryStore(), 105*time.Minute, 20, 30)
	handler := NewHandler(authService, presenceService)
	router := NewRouter(handler, authService)

	request := httptest.NewRequest(http.MethodPost, "/presence", bytes.NewReader([]byte(`{"activity":"SELF_STUDY"}`)))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d: %s", response.Code, response.Body.String())
	}
}

func requestToken(t *testing.T, router http.Handler) string {
	t.Helper()

	request := httptest.NewRequest(http.MethodPost, "/auth/discord/callback", bytes.NewReader([]byte(`{"code":"test-code"}`)))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected auth status 200, got %d: %s", response.Code, response.Body.String())
	}

	var body generated.AuthTokenResponseBody
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.AccessToken == "" {
		t.Fatal("expected access token")
	}

	return body.AccessToken
}
