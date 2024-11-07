package repositories

import (
    "time"
    "errors"
    "context"

    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/mykykh/concerts-api/internal/domain"
)

func Save(db *pgxpool.Pool, concert domain.Concert) error {
    _, err := db.Exec(
        context.Background(),
        "INSERT INTO concerts(title, description, location) VALUES ($1, $2, $3)",
        concert.Title,
        concert.Description,
        concert.Location,
    );

    if err != nil {
        return errors.New("Failed to save concert")
    }

    return nil
}

func GetById(db *pgxpool.Pool, id int64) (*domain.Concert, error) {
    var title, description, location string;
    var createDate, updateDate time.Time;

    err := db.QueryRow(
        context.Background(),
        "SELECT title, description, location, create_date, update_date FROM concerts WHERE concert_id=$1",
        id,
    ).Scan(&title, &description, &location, &createDate, &updateDate);

    if err != nil {
        return nil, errors.New("Failed to get concert")
    }

    return &domain.Concert{
        ID: id,
        Title: title,
        Description: description,
        Location: location,
        CreateDate: createDate,
        UpdateDate: updateDate,
    }, nil
}
