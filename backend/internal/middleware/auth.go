package middleware

import (
	"context"
	"net/http"
	"strings"

	"lms/internal/config"
	"lms/internal/models"
	"lms/internal/utils"
)

type contextKey string

const UserContextKey contextKey = "authenticated_user"

func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := ""

			// 1. Try Authorization header
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			}

			// 2. Try HttpOnly Cookie
			if tokenStr == "" {
				if cookie, err := r.Cookie("lms_token"); err == nil {
					tokenStr = cookie.Value
				}
			}

			// 3. Try query parameter (useful for WebSocket connection handshake)
			if tokenStr == "" {
				tokenStr = r.URL.Query().Get("token")
			}

			if tokenStr != "" {
				claims, err := utils.ValidateJWT(tokenStr, cfg.JWTSecret)
				if err == nil && claims != nil {
					user := &models.User{
						ID:    claims.UserID,
						Email: claims.Email,
						Role:  models.Role(claims.Role),
						Name:  claims.Name,
					}
					ctx := context.WithValue(r.Context(), UserContextKey, user)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user == nil {
			utils.Unauthorized(w, "Authentication required. Please log in.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireRoles(roles ...models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUser(r)
			if user == nil {
				utils.Unauthorized(w, "Authentication required. Please log in.")
				return
			}

			allowed := false
			for _, role := range roles {
				if user.Role == role {
					allowed = true
					break
				}
			}

			if !allowed {
				utils.Forbidden(w, "You do not have permission to access this resource.")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireAdmin(next http.Handler) http.Handler {
	return RequireRoles(models.RoleAdmin)(next)
}

func RequireTeacherOrAdmin(next http.Handler) http.Handler {
	return RequireRoles(models.RoleTeacher, models.RoleAdmin)(next)
}

func GetUser(r *http.Request) *models.User {
	if val := r.Context().Value(UserContextKey); val != nil {
		if user, ok := val.(*models.User); ok {
			return user
		}
	}
	return nil
}
