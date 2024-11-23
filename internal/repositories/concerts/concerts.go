package concerts

import (
    "time"
    "errors"
    "context"

    "github.com/gofrs/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/mykykh/concerts-api/internal/domain"
)

func Save(db *pgxpool.Pool, concert domain.Concert) error {
    _, err := db.Exec(
        context.Background(),
        "INSERT INTO concerts(title, description, location, author_id) VALUES ($1, $2, $3, $4)",
        concert.Title,
        concert.Description,
        concert.Location,
        concert.AuthorID,
    );

    if err != nil {
        return errors.New("Failed to save concert")
    }

    return nil
}

func parseRows(rows pgx.Rows) ([]domain.Concert, error) {
    var concerts []domain.Concert

    for rows.Next() {
        var concert domain.Concert
        err := rows.Scan(
            &concert.ID,
            &concert.Title,
            &concert.Description,
            &concert.Location,
            &concert.AuthorID,
            &concert.CreateDate,
            &concert.UpdateDate,
        )
        if err != nil {
            return []domain.Concert{}, errors.New("Failed to parse rows")
        }
        concerts = append(concerts, concert)
    }

    return concerts, nil
}

func GetLast10(db *pgxpool.Pool) ([]domain.Concert, error) {
    rows, err := db.Query(
        context.Background(),
        "SELECT concert_id, title, description, location, author_id, create_date, update_date FROM concerts ORDER BY concert_id DESC LIMIT 10",
    )

    if err != nil {
        return []domain.Concert{}, errors.New("Failed to select rows")
    }

    return parseRows(rows)
}

func GetLast10BeforeId(db *pgxpool.Pool, concertId int64) ([]domain.Concert, error) {
    rows, err := db.Query(
        context.Background(),
        "SELECT concert_id, title, description, location, author_id, create_date, update_date FROM concerts WHERE concert_id < $1 ORDER BY concert_id DESC LIMIT 10",
        concertId,
    )

    if err != nil {
        return []domain.Concert{}, errors.New("Failed to select rows")
    }

    return parseRows(rows)
}

func GetById(db *pgxpool.Pool, id int64) (*domain.Concert, error) {
    var authorId uuid.UUID;
    var title, description, location string;
    var createDate, updateDate time.Time;

    err := db.QueryRow(
        context.Background(),
        "SELECT title, description, location, author_id, create_date, update_date FROM concerts WHERE concert_id=$1",
        id,
    ).Scan(&title, &description, &location, &authorId, &createDate, &updateDate);

    if err != nil {
        return nil, errors.New("Failed to get concert")
    }

    return &domain.Concert{
        ID: id,
        Title: title,
        Description: description,
        Location: location,
        AuthorID: authorId,
        CreateDate: createDate,
        UpdateDate: updateDate,
    }, nil
}

func Update(db *pgxpool.Pool, concert domain.Concert) error {
    _, err := db.Exec(
        context.Background(),
        "UPDATE concerts SET title=$2, description=$3, location=$4 WHERE concert_id=$1",
        concert.ID,
        concert.Title,
        concert.Description,
        concert.Location,
    )

    if err != nil {
        return errors.New("Failed to update concert")
    }

    return nil
}
