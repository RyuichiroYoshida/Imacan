package server

import (
	"context"
	"time"

	"github.com/RyuichiroYoshida/imacan/backend/internal/auth"
	"github.com/RyuichiroYoshida/imacan/backend/internal/generated"
	"github.com/RyuichiroYoshida/imacan/backend/internal/presence"
)

type Handler struct {
	auth     *auth.Service
	presence *presence.Service
	now      func() time.Time
}

func NewHandler(authService *auth.Service, presenceService *presence.Service) *Handler {
	return &Handler{
		auth:     authService,
		presence: presenceService,
		now:      time.Now,
	}
}

func (h *Handler) AuthDiscordCallback(ctx context.Context, request generated.AuthDiscordCallbackRequestObject) (generated.AuthDiscordCallbackResponseObject, error) {
	if request.Body == nil || request.Body.Code == "" {
		return generated.AuthDiscordCallback400JSONResponse(errorBody("INVALID_REQUEST", "code is required")), nil
	}

	userID := auth.UserIDFromDiscordCode(request.Body.Code)
	token, expiresIn, err := h.auth.Issue(userID, h.now())
	if err != nil {
		return generated.AuthDiscordCallback500JSONResponse(errorBody("TOKEN_ISSUE_FAILED", "failed to issue token")), nil
	}

	return generated.AuthDiscordCallback200JSONResponse{
		AccessToken: token,
		TokenType:   generated.Bearer,
		ExpiresIn:   expiresIn,
	}, nil
}

func (h *Handler) PresenceUpdatePresence(ctx context.Context, request generated.PresenceUpdatePresenceRequestObject) (generated.PresenceUpdatePresenceResponseObject, error) {
	userID, ok := userIDFromContext(ctx)
	if !ok {
		return generated.PresenceUpdatePresence401JSONResponse(errorBody("UNAUTHORIZED", "valid bearer token is required")), nil
	}
	if request.Body == nil {
		return generated.PresenceUpdatePresence400JSONResponse(errorBody("INVALID_REQUEST", "request body is required")), nil
	}

	record, hasExpiry, err := h.presence.Update(userID, request.Body.Activity, h.now())
	if err != nil {
		return generated.PresenceUpdatePresence400JSONResponse(errorBody("INVALID_ACTIVITY", "activity must be CLASS, SELF_STUDY, or OUT")), nil
	}

	response := generated.PresenceUpdatePresence200JSONResponse{
		Activity:  record.Activity,
		UpdatedAt: record.UpdatedAt,
	}
	if hasExpiry {
		response.ExpiresAt = &record.ExpiresAt
	}

	return response, nil
}

func (h *Handler) PresenceGetPresenceSummary(ctx context.Context, request generated.PresenceGetPresenceSummaryRequestObject) (generated.PresenceGetPresenceSummaryResponseObject, error) {
	summary := h.presence.Summary(h.now())

	return generated.PresenceGetPresenceSummary200JSONResponse{
		Total:     summary.Total,
		Class:     summary.Class,
		SelfStudy: summary.SelfStudy,
	}, nil
}

func errorBody(code, message string) generated.ErrorResponseBody {
	return generated.ErrorResponseBody{
		Code:    code,
		Message: message,
	}
}
