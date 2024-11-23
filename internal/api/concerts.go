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

// @Summary Returns list of concerts
// @Description Returns list of 10 last concerts in db oredered by concert_id.
// @Description Use last_id query parameter to select concerts before this last_id
// @Tags Concerts
// @Param last_id query integer false "id before which to get last concerts"
// @Router /concerts [get]
func (rs ConcertsResource) GetAll(w http.ResponseWriter, r *http.Request) {
    lastId := r.URL.Query().Get("last_id")

    var concerts []domain.Concert
    var err error

    if lastId == "" {
        concerts, err = concertsRepository.GetLast10(rs.db)
    } else {
        lastId, err := strconv.ParseInt(lastId, 10, 64)
        if err != nil {
            http.Error(w, "Failed to parse last_id", http.StatusBadRequest)
            return
        }
        concerts, err = concertsRepository.GetLast10BeforeId(rs.db, lastId)
    }

    if err != nil {
        http.Error(w, "Failed to find concerts", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(concerts)
}

// @Summary Creates new concert
// @Tags Concerts
// @Param concert body domain.Concert true "Concert to create"
// @Security BearerAuth
// @Router /concerts [post]
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

    concert.AuthorID = claims.ID;

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

// @Summary Returns concert info
// @Tags Concerts
// @Param id path integer true "Concert id"
// @Router /concerts/{id} [get]
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

// @Summary Updates concert info
// @Tags Concerts
// @Param id path integer true "Concert id"
// @Param concert body domain.Concert true "Updated concert info"
// @Security BearerAuth
// @Router /concerts/{id} [put]
func (rs ConcertsResource) Update(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value("claims").(auth.Claims)
    if !ok {
        http.Error(w, "No claims found", http.StatusUnauthorized)
        return
    }

    if !claims.HasResourceAccessRole("concerts-api", "concerts-update-any") {
        if !claims.HasResourceAccessRole("concerts-api", "concerts-update") {
            http.Error(w, "No concerts-update role found", http.StatusUnauthorized)
            return
        }
    }

    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, http.StatusText(400), 400)
        return
    }

    concert, err := concertsRepository.GetById(rs.db, id)
    if err != nil {
        http.Error(w, "Concert with id not found", http.StatusNotFound)
        return
    }

    if !claims.HasResourceAccessRole("concerts-api", "concerts-update-any") {
        if concert.AuthorID != claims.ID {
            http.Error(w, "Unauthorized to update concert with id", http.StatusUnauthorized)
        }
    }

    var updatedConcert domain.Concert
    err = json.NewDecoder(r.Body).Decode(&updatedConcert)

    if err != nil {
        http.Error(w, http.StatusText(400), 400)
        return
    }

    updatedConcert.ID = id

    err = concertsRepository.Update(rs.db, updatedConcert);

    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }
}
