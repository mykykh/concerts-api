package api

import (
    "net/http"

    "github.com/go-chi/chi/v5"
)

type ConcertsResource struct {}

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
    w.Write([]byte("create"))
}

func (rs ConcertsResource) Get(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("get"))
}

func (rs ConcertsResource) Update(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("update"))
}
