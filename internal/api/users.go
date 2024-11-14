package api

import (
    "net/http"
    "encoding/json"

    "github.com/gofrs/uuid"
    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"

    usersRepository "github.com/mykykh/concerts-api/internal/repositories/users"
)

type UsersResource struct {
    db *pgxpool.Pool
}

func (rs UsersResource) Routes() chi.Router {
    r := chi.NewRouter()

    r.Get("/", rs.GetAll)

    r.Route("/{id}", func (r chi.Router) {
        r.Get("/", rs.Get)
    })

    return r
}

func (rs UsersResource) GetAll(w http.ResponseWriter, r *http.Request) {
    users, err := usersRepository.GetAll(rs.db)

    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    json.NewEncoder(w).Encode(users)
}

func (rs UsersResource) Get(w http.ResponseWriter, r *http.Request) {
    id, err := uuid.FromString(chi.URLParam(r, "id"))
    if err != nil {
        http.Error(w, http.StatusText(400), 400)
        return
    }

    user, err := usersRepository.GetById(rs.db, id)

    if err != nil {
        http.Error(w, http.StatusText(404), 404)
        return
    }

    json.NewEncoder(w).Encode(user)
}
