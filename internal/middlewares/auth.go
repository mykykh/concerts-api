package middlewares

import (
    "os"
    "fmt"
    "strings"
    "context"
    "net/http"

    "golang.org/x/oauth2"
    "github.com/coreos/go-oidc/v3/oidc"

    "github.com/mykykh/concerts-api/internal/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
            return
        }

        // Expecting "Bearer <token>"
        token := strings.TrimPrefix(authHeader, "Bearer ")
        if token == authHeader {
            http.Error(w, "Bearer token missing", http.StatusUnauthorized)
            return
        }

        provider, err := oidc.NewProvider(oauth2.NoContext, os.Getenv("OAUTH_SERVER_HOSTNAME"))
        if err != nil {
            return
        }
        verifier := provider.Verifier(&oidc.Config{ClientID: os.Getenv("OAUTH_CLIENT_ID")})

        // Verify the ID token
        idToken, err := verifier.Verify(oauth2.NoContext, token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        var claims auth.Claims
        if err := idToken.Claims(&claims); err != nil {
            // handle error
            http.Error(w, "Failed to parse claims", http.StatusUnauthorized)
            return
        }

        // Token is valid, pass the request to the next handler
        r = r.WithContext(context.WithValue(r.Context(), "claims", claims))
        next.ServeHTTP(w, r)
    })
}
