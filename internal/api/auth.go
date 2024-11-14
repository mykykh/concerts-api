package api

import (
    "os"
    "fmt"
    "errors"
    "net/http"
    "encoding/json"

    "golang.org/x/oauth2"
    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/coreos/go-oidc/v3/oidc"

    "github.com/mykykh/concerts-api/internal/auth"
    "github.com/mykykh/concerts-api/internal/domain"
    usersRepository "github.com/mykykh/concerts-api/internal/repositories/users"
)

type AuthResource struct {
    db *pgxpool.Pool
}

func (rs AuthResource) Routes() chi.Router {
    r := chi.NewRouter()

    r.Get("/login", rs.handleRedirect)
    r.Get("/callback", rs.handleOAuth2Callback)

    return r
}

func initOauth() (*oauth2.Config, error) {
    provider, err := oidc.NewProvider(oauth2.NoContext, os.Getenv("OAUTH_SERVER_HOSTNAME"))
    if err != nil {
        fmt.Println(err)
        return nil, errors.New("Failed to create oidc provider")
    }

    // Configure an OpenID Connect aware OAuth2 client.
    oauth2Config := oauth2.Config{
        ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
        ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
        RedirectURL: "http://localhost:8080/auth/callback",

        // Discovery returns the OAuth2 endpoints.
        Endpoint: provider.Endpoint(),

        // "openid" is a required scope for OpenID Connect flows.
        Scopes: []string{oidc.ScopeOpenID, "email"},
    }

    return &oauth2Config, nil
}

func (rs AuthResource) handleRedirect(w http.ResponseWriter, r *http.Request) {
    state := "random_state"
    oauth2Config, err := initOauth()
    if err != nil {
        return
    }
    http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
}

func saveTokenUserToDB(db *pgxpool.Pool, token *oauth2.Token) error {

    provider, err := oidc.NewProvider(oauth2.NoContext, os.Getenv("OAUTH_SERVER_HOSTNAME"))
    if err != nil {
        return errors.New("Failed to connect to provider")
    }
    verifier := provider.Verifier(&oidc.Config{ClientID: os.Getenv("OAUTH_CLIENT_ID")})

    idToken, err := verifier.Verify(oauth2.NoContext, token.AccessToken)
    if err != nil {
        return errors.New("Failed to verify token")
    }

    claims, err := auth.TokenToClaims(idToken)
    if err != nil {
        fmt.Println(err)
        return errors.New("Failed to parse claims")
    }

    user := domain.ClaimsToUser(*claims)

    err = usersRepository.SaveOrUpdate(db, user)

    if err != nil {
        fmt.Println(err)
        return errors.New("Failed to save user")
    }
    return nil
}

func (rs AuthResource) handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
    oauth2Config, err := initOauth()
    if err != nil {
        return
    }
    oauth2Token, err := oauth2Config.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
    if err != nil {
        fmt.Println(err)
        return
    }

    err = saveTokenUserToDB(rs.db, oauth2Token)
    if err != nil {
        fmt.Println(err)
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "access-token": oauth2Token.AccessToken,
    })
}
