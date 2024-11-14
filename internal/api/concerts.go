package api

import (
    "strconv"
    "net/http"
    "encoding/json"

    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/mykykh/concerts-api/internal/auth"
    "github.com/mykykh/concerts-api/internal/domain"
    concertsRepository "github.com/mykykh/concerts-api/internal/repositories/concerts"
    "github.com/mykykh/concerts-api/internal/middlewares"
)

type ConcertsResource struct {
    db *pgxpool.Pool
}

func (rs ConcertsResource) Routes() chi.Router {
    r := chi.NewRouter()

    r.Get("/", rs.GetAll)
    r.With(middlewares.AuthMiddleware).Post("/", rs.Create)

    r.Route("/{id}", func (r chi.Router) {
        r.Get("/", rs.Get)
        r.With(middlewares.AuthMiddleware).Put("/", rs.Update)
    })

    return r
}

func (rs ConcertsResource) GetAll(w http.ResponseWriter, r *http.Request) {
    concerts, err := concertsRepository.GetAll(rs.db)

    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    json.NewEncoder(w).Encode(concerts)
}

func (rs ConcertsResource) Create(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value("claims").(auth.Claims)
    if !ok {
        http.Error(w, "No claims found", http.StatusUnauthorized)
        return
    }

    if !claims.HasResourceAccessRole("concerts-api", "concerts-create") {
        http.Error(w, "No concerts-create role found", http.StatusUnauthorized)
        return
    }

    var concert domain.Concert;

    err := json.NewDecoder(r.Body).Decode(&concert);

    if err != nil {
        http.Error(w, http.StatusText(400), 400)
        return
    }

    err = concertsRepository.Save(rs.db, concert);

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

    concert, err := concertsRepository.GetById(rs.db, id)

    if err != nil {
        http.Error(w, http.StatusText(404), 404)
        return
    }

    json.NewEncoder(w).Encode(concert)
}

func (rs ConcertsResource) Update(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value("claims").(auth.Claims)
    if !ok {
        http.Error(w, "No claims found", http.StatusUnauthorized)
        return
    }

    if !claims.HasResourceAccessRole("concerts-api", "concerts-update") {
        http.Error(w, "No concerts-update role found", http.StatusUnauthorized)
        return
    }
    var concert domain.Concert

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

    err = concertsRepository.Update(rs.db, concert);

    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }
}
