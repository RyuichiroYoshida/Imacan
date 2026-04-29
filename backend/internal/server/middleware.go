package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/RyuichiroYoshida/imacan/backend/internal/auth"
	"github.com/RyuichiroYoshida/imacan/backend/internal/generated"
)

func AuthMiddleware(authService *auth.Service) generated.StrictMiddlewareFunc {
	return func(next generated.StrictHandlerFunc, operationID string) generated.StrictHandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
			if operationID != "PresenceUpdatePresence" {
				return next(ctx, w, r, request)
			}

			token := bearerToken(r.Header.Get("Authorization"))
			if token == "" {
				return generated.PresenceUpdatePresence401JSONResponse(errorBody("UNAUTHORIZED", "valid bearer token is required")), nil
			}

			claims, err := authService.Verify(token, time.Now())
			if err != nil {
				return generated.PresenceUpdatePresence401JSONResponse(errorBody("UNAUTHORIZED", "valid bearer token is required")), nil
			}

			return next(withUserID(ctx, claims.UserID), w, r, request)
		}
	}
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}
