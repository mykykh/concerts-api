package api

import (
    "strconv"
    "net/http"
    "encoding/json"

    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/mykykh/concerts-api/internal/domain"
    "github.com/mykykh/concerts-api/internal/repositories"
    "github.com/mykykh/concerts-api/internal/middlewares"
)

type ConcertsResource struct {
    db *pgxpool.Pool
}

func (rs ConcertsResource) Routes() chi.Router {
    r := chi.NewRouter()

    r.Get("/", rs.GetAll)
    r.Post("/", rs.Create)

    r.Route("/{id}", func (r chi.Router) {
        r.Use(middlewares.AuthMiddleware)
        r.Get("/", rs.Get)
        r.Put("/", rs.Update)
    })

    return r
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
