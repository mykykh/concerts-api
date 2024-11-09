package api

import (
    "os"
    "fmt"
    "errors"
    "strconv"
    "net/http"
    "encoding/json"

    "golang.org/x/oauth2"
    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/coreos/go-oidc/v3/oidc"

    "github.com/mykykh/concerts-api/internal/domain"
    "github.com/mykykh/concerts-api/internal/repositories"
)

type ConcertsResource struct {
    db *pgxpool.Pool
}

func (rs ConcertsResource) Routes() chi.Router {
    r := chi.NewRouter()

    r.Get("/", rs.GetAll)
    r.Post("/", rs.Create)

    r.Get("/login", rs.handleRedirect)
    r.Get("/callback", rs.handleOAuth2Callback)

    r.Route("/{id}", func (r chi.Router) {
        r.Get("/", rs.Get)
        r.Put("/", rs.Update)
    })

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
        RedirectURL: "http://localhost:8080/concerts/callback",

        // Discovery returns the OAuth2 endpoints.
        Endpoint: provider.Endpoint(),

        // "openid" is a required scope for OpenID Connect flows.
        Scopes: []string{oidc.ScopeOpenID, "email"},
    }

    return &oauth2Config, nil
}

func (rs ConcertsResource) handleRedirect(w http.ResponseWriter, r *http.Request) {
    state := "random_state"
    oauth2Config, err := initOauth()
    if err != nil {
        fmt.Println(err)
        return
        // handle error
    }
    http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
}


func (rs ConcertsResource) handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
    oauth2Config, err := initOauth()
    oauth2Token, err := oauth2Config.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
    if err != nil {
        fmt.Println(err)
        return
        // handle error
    }

    // Extract the ID Token from OAuth2 token.
    rawIDToken, ok := oauth2Token.Extra("id_token").(string)
    if !ok {
        fmt.Println(err)
        return
        // handle missing token
    }

    provider, err := oidc.NewProvider(oauth2.NoContext, os.Getenv("OAUTH_SERVER_HOSTNAME"))
    if err != nil {
        fmt.Println(err)
        return
    }
    verifier := provider.Verifier(&oidc.Config{ClientID: os.Getenv("OAUTH_CLIENT_ID")})
    idToken, err := verifier.Verify(oauth2.NoContext, rawIDToken)
    if err != nil {
        fmt.Println(err)
        return
        // handle error
    }

    // Extract custom claims
    var claims struct {
        Email    string `json:"email"`
        Verified bool   `json:"email_verified"`
    }
    if err := idToken.Claims(&claims); err != nil {
        // handle error
    }

    fmt.Println(claims)
}

func (rs ConcertsResource) GetAll(w http.ResponseWriter, r *http.Request) {
    concerts, err := repositories.GetAll(rs.db)

    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    json.NewEncoder(w).Encode(concerts)
}

func (rs ConcertsResource) Create(w http.ResponseWriter, r *http.Request) {
    var concert domain.Concert;

    err := json.NewDecoder(r.Body).Decode(&concert);

    if err != nil {
        http.Error(w, http.StatusText(400), 400)
        return
    }

    err = repositories.Save(rs.db, concert);

    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }
}

func (rs ConcertsResource) Get(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, http.StatusText(400), 400)
        return
    }

    concert, err := repositories.GetById(rs.db, id)

    if err != nil {
        http.Error(w, http.StatusText(404), 404)
        return
    }

    json.NewEncoder(w).Encode(concert)
}

func (rs ConcertsResource) Update(w http.ResponseWriter, r *http.Request) {
    var concert domain.Concert;

    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, http.StatusText(400), 400)
        return
    }

    err = json.NewDecoder(r.Body).Decode(&concert)

    if err != nil {
        http.Error(w, http.StatusText(400), 400)
        return
    }

    concert.ID = id

    err = repositories.Update(rs.db, concert);

    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }
}
