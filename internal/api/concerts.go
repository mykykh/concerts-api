package api

import (
    "time"
    "strconv"
    "net/http"
    "encoding/json"

    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"

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

    r.Route("/{id}", func (r chi.Router) {
        r.Get("/", rs.Get)
        r.Put("/", rs.Update)
    })

    return r
}

func (rs ConcertsResource) GetAll(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("getall"))
}

func (rs ConcertsResource) Create(w http.ResponseWriter, r *http.Request) {
    concert := domain.Concert {
        ID: 0,
        Title: "test",
        Description: "",
        Location: "test",
        CreateDate: time.Now(),
        UpdateDate: time.Now(),
    };
    repositories.Save(rs.db, concert);
}

func (rs ConcertsResource) Get(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        return
    }
    concert := repositories.GetById(rs.db, id)
    json.NewEncoder(w).Encode(concert)
}

func (rs ConcertsResource) Update(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("update"))
}
