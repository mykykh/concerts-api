package tickets

import (
    "errors"
    "context"

    "github.com/gofrs/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/mykykh/concerts-api/internal/domain"
)

func Save(db *pgxpool.Pool, ticket domain.Ticket) error {
    _, err := db.Exec(
        context.Background(),
        "INSERT INTO tickets(verification_token, concert_id, user_id) VALUES ($1, $2, $3)",
        ticket.VerificationToken,
        ticket.ConcertID,
        ticket.UserID,
    );

    if err != nil {
        return errors.New("Failed to save ticket")
    }

    return nil
}

func parseRows(rows pgx.Rows) ([]domain.Ticket, error) {
    var tickets []domain.Ticket

    for rows.Next() {
        var ticket domain.Ticket
        err := rows.Scan(
            &ticket.ID,
            &ticket.VerificationToken,
            &ticket.ConcertID,
            &ticket.UserID,
            &ticket.CreateDate,
            &ticket.UpdateDate,
        )
        if err != nil {
            return []domain.Ticket{}, errors.New("Failed to parse rows")
        }
        tickets = append(tickets, ticket)
    }

    return tickets, nil
}

func GetLast10(db *pgxpool.Pool) ([]domain.Ticket, error) {
    rows, err := db.Query(
        context.Background(),
        "SELECT ticket_id, verification_token, concert_id, user_id, create_date, update_date FROM tickets ORDER BY ticket_id DESC LIMIT 10",
    )

    if err != nil {
        return []domain.Ticket{}, errors.New("Failed to select rows")
    }

    return parseRows(rows)
}

func GetLast10BeforeId(db *pgxpool.Pool, ticketId int64) ([]domain.Ticket, error) {
    rows, err := db.Query(
        context.Background(),
        "SELECT ticket_id, verification_token, concert_id, user_id, create_date, update_date FROM tickets WHERE ticket_id < $1 ORDER BY ticket_id DESC LIMIT 10",
        ticketId,
    )

    if err != nil {
        return []domain.Ticket{}, errors.New("Failed to select rows")
    }

    return parseRows(rows)
}

func GetOwnLast10(db *pgxpool.Pool, userId uuid.UUID) ([]domain.Ticket, error) {
    rows, err := db.Query(
        context.Background(),
        "SELECT ticket_id, verification_token, concert_id, user_id, create_date, update_date FROM tickets WHERE user_id = $1 ORDER BY ticket_id DESC LIMIT 10",
        userId,
    )

    if err != nil {
        return []domain.Ticket{}, errors.New("Failed to select rows")
    }

    return parseRows(rows)
}

func GetOwnLast10BeforeId(db *pgxpool.Pool, ticketId int64, userId uuid.UUID) ([]domain.Ticket, error) {
    rows, err := db.Query(
        context.Background(),
        "SELECT ticket_id, verification_token, concert_id, user_id, create_date, update_date FROM tickets WHERE ticket_id < $1 AND user_id = $2 ORDER BY ticket_id DESC LIMIT 10",
        ticketId,
        userId,
    )

    if err != nil {
        return []domain.Ticket{}, errors.New("Failed to select rows")
    }

    return parseRows(rows)
}

func GetById(db *pgxpool.Pool, id int64) (*domain.Ticket, error) {
    var ticket domain.Ticket

    err := db.QueryRow(
        context.Background(),
        "SELECT ticket_id, verification_token, concert_id, user_id, create_date, update_date FROM tickets WHERE ticket_id=$1",
        id,
    ).Scan(
        &ticket.ID,
        &ticket.VerificationToken,
        &ticket.ConcertID,
        &ticket.UserID,
        &ticket.CreateDate,
        &ticket.UpdateDate,
    );

    if err != nil {
        return nil, errors.New("Failed to get ticket")
    }

    return &ticket, nil
}

func Update(db *pgxpool.Pool, ticket domain.Ticket) error {
    _, err := db.Exec(
        context.Background(),
        "UPDATE tickets SET verification_token=$2 WHERE ticket_id=$1",
        ticket.ID,
        ticket.VerificationToken,
    )

    if err != nil {
        return errors.New("Failed to update ticket")
    }

    return nil
}
