package concerts

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

func GetAll(db *pgxpool.Pool) ([]domain.Concert, error) {
    var concerts []domain.Concert

    rows, err := db.Query(
        context.Background(),
        "SELECT concert_id, title, description, location, create_date, update_date FROM concerts",
    )

    if err != nil {
        return []domain.Concert{}, errors.New("Failed to select rows")
    }

    for rows.Next() {
        var concert domain.Concert
        err := rows.Scan(&concert.ID, &concert.Title, &concert.Description, &concert.Location, &concert.CreateDate, &concert.UpdateDate)
        if err != nil {
            return []domain.Concert{}, errors.New("Failed to parse rows")
        }
        concerts = append(concerts, concert)
    }

    return concerts, nil
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
