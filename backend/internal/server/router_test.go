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
	router, _ := newTestRouter(t)

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

func TestCurrentPresence(t *testing.T) {
	router, _ := newTestRouter(t)
	token := requestToken(t, router)

	updateBody := []byte(`{"activity":"CLASS"}`)
	updateRequest := httptest.NewRequest(http.MethodPost, "/presence", bytes.NewReader(updateBody))
	updateRequest.Header.Set("Authorization", "Bearer "+token)
	updateRequest.Header.Set("Content-Type", "application/json")
	updateResponse := httptest.NewRecorder()
	router.ServeHTTP(updateResponse, updateRequest)
	if updateResponse.Code != http.StatusOK {
		t.Fatalf("expected update status 200, got %d: %s", updateResponse.Code, updateResponse.Body.String())
	}

	request := httptest.NewRequest(http.MethodGet, "/presence/me", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("expected current presence status 200, got %d: %s", response.Code, response.Body.String())
	}

	var body generated.CurrentPresenceResponseBody
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if !body.Active || body.Activity == nil || *body.Activity != generated.CLASS || body.ExpiresAt == nil {
		t.Fatalf("unexpected current presence: %+v", body)
	}
}

func TestCurrentPresenceReturnsInactiveWhenExpired(t *testing.T) {
	router, handler := newTestRouter(t)
	token := requestToken(t, router)

	handler.now = func() time.Time {
		return time.Date(2026, 4, 26, 20, 29, 0, 0, time.UTC)
	}
	updateRequest := httptest.NewRequest(http.MethodPost, "/presence", bytes.NewReader([]byte(`{"activity":"SELF_STUDY"}`)))
	updateRequest.Header.Set("Authorization", "Bearer "+token)
	updateRequest.Header.Set("Content-Type", "application/json")
	updateResponse := httptest.NewRecorder()
	router.ServeHTTP(updateResponse, updateRequest)
	if updateResponse.Code != http.StatusOK {
		t.Fatalf("expected update status 200, got %d: %s", updateResponse.Code, updateResponse.Body.String())
	}

	handler.now = func() time.Time {
		return time.Date(2026, 4, 26, 20, 31, 0, 0, time.UTC)
	}
	request := httptest.NewRequest(http.MethodGet, "/presence/me", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("expected current presence status 200, got %d: %s", response.Code, response.Body.String())
	}

	var body generated.CurrentPresenceResponseBody
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Active {
		t.Fatalf("expected inactive presence, got %+v", body)
	}
}

func TestPresenceRequiresBearerToken(t *testing.T) {
	router, _ := newTestRouter(t)

	request := httptest.NewRequest(http.MethodPost, "/presence", bytes.NewReader([]byte(`{"activity":"SELF_STUDY"}`)))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d: %s", response.Code, response.Body.String())
	}
}

func TestCurrentPresenceRequiresBearerToken(t *testing.T) {
	router, _ := newTestRouter(t)

	request := httptest.NewRequest(http.MethodGet, "/presence/me", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d: %s", response.Code, response.Body.String())
	}
}

func newTestRouter(t *testing.T) (http.Handler, *Handler) {
	t.Helper()

	discord := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/token":
			if err := r.ParseForm(); err != nil {
				t.Fatalf("parse discord token form: %v", err)
			}
			if r.Form.Get("code") != "test-code" || r.Form.Get("grant_type") != "authorization_code" {
				t.Fatalf("unexpected token form: %v", r.Form)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"discord-access-token"}`))
		case "/users/@me":
			if r.Header.Get("Authorization") != "Bearer discord-access-token" {
				t.Fatalf("unexpected discord authorization header: %s", r.Header.Get("Authorization"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"1234567890"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(discord.Close)

	authService := auth.NewService("test-secret", time.Hour)
	authService.ConfigureDiscord(auth.DiscordConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURI:  "http://localhost:3000/auth/callback",
		TokenURL:     discord.URL + "/oauth2/token",
		UserURL:      discord.URL + "/users/@me",
	}, discord.Client())
	presenceService := presence.NewService(presence.NewMemoryStore(), 105*time.Minute, 20, 30)
	handler := NewHandler(authService, presenceService)

	return NewRouter(handler, authService), handler
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
