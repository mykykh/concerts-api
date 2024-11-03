package main

import (
    "net/http"

    "github.com/mykykh/concerts-api/internal/api"
    "github.com/go-chi/chi/v5"
)

func main() {
    r := chi.NewRouter()

    r.Mount("/concerts", api.ConcertsResource{}.Routes())

    http.ListenAndServe(":8080", r)
}
