package middleware

import (
	"context"
	"net/http"

	"github.com/anden007/af_dp_clean_core/part"
	"github.com/anden007/af_dp_clean_core/pkg"

	"github.com/spf13/viper"
)

// Middleware decodes the share session cookie and packs the session into context
func NewAuthMiddleware(db part.IDataBase, jwt part.IJwtService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtHeader := viper.GetString("jwt.header")
		needAuth := viper.GetBool("gql.need_auth")
		if needAuth {
			if tokenCookie, err := r.Cookie(jwtHeader); err == nil {
				if token, vcErr := jwt.VerifyCookieToken(tokenCookie.Value); vcErr == nil {
					var claims pkg.BaseUserInfoClaims
					if cErr := token.Claims(&claims); cErr == nil {
						ctx := context.WithValue(r.Context(), "jwt.claims", claims.BaseUserInfo)
						r = r.WithContext(ctx)
						next.ServeHTTP(w, r)
					}
				}
				return
			}
		} else {
			// 允许未授权访问
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, `{"success": false,"message": "UnAuth Request"}`, http.StatusForbidden)
		return
	})
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *pkg.BaseUserInfo {
	raw, _ := ctx.Value("jwt.claims").(*pkg.BaseUserInfo)
	return raw
}
