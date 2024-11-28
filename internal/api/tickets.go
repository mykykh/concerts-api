package api

import (
    "strconv"
    "net/http"
    "encoding/json"

    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/mykykh/concerts-api/internal/auth"
    "github.com/mykykh/concerts-api/internal/domain"
    ticketsRepository "github.com/mykykh/concerts-api/internal/repositories/tickets"
    "github.com/mykykh/concerts-api/internal/middlewares"
)

type TicketsResource struct {
    db *pgxpool.Pool
}

func (rs TicketsResource) Routes() chi.Router {
    r := chi.NewRouter()

    r.With(middlewares.AuthMiddleware).Get("/", rs.GetAll)
    r.With(middlewares.AuthMiddleware).Get("/own", rs.GetAllOwn)
    r.With(middlewares.AuthMiddleware).Post("/", rs.Create)

    r.Route("/{id}", func (r chi.Router) {
        r.With(middlewares.AuthMiddleware).Get("/", rs.Get)
        r.With(middlewares.AuthMiddleware).Put("/", rs.Update)
    })

    return r
}

// @Summary Returns list of tickets
// @Description Returns list of 10 last tickets in db oredered by ticket_id.
// @Description Use last_id query parameter to select tickets before this last_id
// @Tags Tickets
// @Param last_id query integer false "id before which to get last tickets"
// @Security BearerAuth
// @Router /tickets [get]
func (rs TicketsResource) GetAll(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value("claims").(auth.Claims)
    if !ok {
        http.Error(w, "No claims found", http.StatusUnauthorized)
        return
    }

    if !claims.HasResourceAccessRole("concerts-api", "read:tickets") {
        http.Error(w, "Permission to read tickets not found", http.StatusUnauthorized)
        return
    }
    lastId := r.URL.Query().Get("last_id")

    var tickets []domain.Ticket
    var err error

    if lastId == "" {
        tickets, err = ticketsRepository.GetLast10(rs.db)
    } else {
        lastId, err := strconv.ParseInt(lastId, 10, 64)
        if err != nil {
            http.Error(w, "Failed to parse last_id", http.StatusBadRequest)
            return
        }
        tickets, err = ticketsRepository.GetLast10BeforeId(rs.db, lastId)
    }

    if err != nil {
        http.Error(w, "Failed to find tickets", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(tickets)
}

// @Summary Returns list of tickets that user has
// @Description Returns list of 10 last tickets that the user has in db oredered by ticket_id.
// @Description Use last_id query parameter to select tickets before this last_id
// @Tags Tickets
// @Param last_id query integer false "id before which to get last tickets"
// @Security BearerAuth
// @Router /tickets/own [get]
func (rs TicketsResource) GetAllOwn(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value("claims").(auth.Claims)
    if !ok {
        http.Error(w, "No claims found", http.StatusUnauthorized)
        return
    }

    if !claims.HasResourceAccessRole("concerts-api", "read:tickets") {
        if !claims.HasResourceAccessRole("concerts-api", "read:ownTickets") {
            http.Error(w, "Permission to read own tickets not found", http.StatusUnauthorized)
            return
        }
    }
    lastId := r.URL.Query().Get("last_id")

    var tickets []domain.Ticket
    var err error

    if lastId == "" {
        tickets, err = ticketsRepository.GetOwnLast10(rs.db, claims.ID)
    } else {
        lastId, err := strconv.ParseInt(lastId, 10, 64)
        if err != nil {
            http.Error(w, "Failed to parse last_id", http.StatusBadRequest)
            return
        }
        tickets, err = ticketsRepository.GetOwnLast10BeforeId(rs.db, lastId, claims.ID)
    }

    if err != nil {
        http.Error(w, "Failed to find tickets", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(tickets)
}

// @Summary Creates new ticket
// @Tags Tickets
// @Param ticket body domain.Ticket true "Ticket to create"
// @Security BearerAuth
// @Router /tickets [post]
func (rs TicketsResource) Create(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value("claims").(auth.Claims)
    if !ok {
        http.Error(w, "No claims found", http.StatusUnauthorized)
        return
    }

    var ticket domain.Ticket;

    err := json.NewDecoder(r.Body).Decode(&ticket);

    if err != nil {
        http.Error(w, "Failed to parse ticket", http.StatusBadRequest)
        return
    }

    ticket.UserID = claims.ID

    err = ticketsRepository.Save(rs.db, ticket);

    if err != nil {
        http.Error(w, "Failed to save ticket", http.StatusInternalServerError)
        return
    }
}

// @Summary Returns ticket info
// @Tags Tickets
// @Param id path integer true "Ticket id"
// @Security BearerAuth
// @Router /tickets/{id} [get]
func (rs TicketsResource) Get(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, "Failed to parse id", http.StatusBadRequest)
        return
    }

    claims, ok := r.Context().Value("claims").(auth.Claims)
    if !ok {
        http.Error(w, "No claims found", http.StatusUnauthorized)
        return
    }

    if claims.HasResourceAccessRole("concerts-api", "read:tickets") {
        ticket, err := ticketsRepository.GetById(rs.db, id)
        if err != nil {
            http.Error(w, http.StatusText(404), 404)
            return
        }

        json.NewEncoder(w).Encode(ticket)
    } else if claims.HasResourceAccessRole("concerts-api", "read:ownTickets") {
        ticket, err := ticketsRepository.GetById(rs.db, id)

        if err != nil {
            http.Error(w, "Ticket not found", http.StatusNotFound)
            return
        }
        if ticket.UserID != claims.ID {
            http.Error(w, "Permission to view ticket not found", http.StatusUnauthorized)
            return
        }

        json.NewEncoder(w).Encode(ticket)
    } else {
        http.Error(w, "Permission to view ticket not found", http.StatusUnauthorized)
    }
}

// @Summary Updates ticket info
// @Tags Tickets
// @Param id path integer true "Ticket id"
// @Param ticket body domain.Ticket true "Updated ticket info"
// @Security BearerAuth
// @Router /tickets/{id} [put]
func (rs TicketsResource) Update(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value("claims").(auth.Claims)
    if !ok {
        http.Error(w, "No claims found", http.StatusUnauthorized)
        return
    }

    id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
    if err != nil {
        http.Error(w, "Failed to parse id", http.StatusBadRequest)
        return
    }

    if claims.HasResourceAccessRole("concerts-api", "update:tickets") {
    } else if claims.HasResourceAccessRole("concerts-api", "update:ownTickets") {
        ticket, err := ticketsRepository.GetById(rs.db, id)
        if err != nil {
            http.Error(w, "Ticket not found", http.StatusNotFound)
            return
        }

        if ticket.UserID != claims.ID {
            http.Error(w, "Permsision to update ticket not found", http.StatusUnauthorized)
            return
        }
    } else {
        http.Error(w, "Permsision to update ticket not found", http.StatusUnauthorized)
        return
    }

    var updatedTicket domain.Ticket
    err = json.NewDecoder(r.Body).Decode(&updatedTicket)

    if err != nil {
        http.Error(w, "Failed to parse ticket", http.StatusBadRequest)
        return
    }

    updatedTicket.ID = id

    err = ticketsRepository.Update(rs.db, updatedTicket);

    if err != nil {
        http.Error(w, "Failed to save ticket to db", http.StatusInternalServerError)
        return
    }
}
