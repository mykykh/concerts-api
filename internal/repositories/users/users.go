package users

import (
    "time"
    "errors"
    "context"

    "github.com/gofrs/uuid"
    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/mykykh/concerts-api/internal/domain"
)

func SaveOrUpdate(db *pgxpool.Pool, user domain.User) error {
    dbUser, err := GetById(db, user.ID)

    if (dbUser == nil) {
        err = Save(db, user)
        return err
    }

    err = Update(db, user)
    return err
}


func Save(db *pgxpool.Pool, user domain.User) error {
    _, err := db.Exec(
        context.Background(),
        "INSERT INTO users(user_id, username, email, full_name) VALUES ($1, $2, $3, $4)",
        user.ID,
        user.Username,
        user.Email,
        user.FullName,
    );

    if err != nil {
        return errors.New("Failed to save user")
    }

    return nil
}

func GetAll(db *pgxpool.Pool) ([]domain.User, error) {
    var users []domain.User

    rows, err := db.Query(
        context.Background(),
        "SELECT user_id, username, email, full_name, create_date, update_date FROM users",
    )

    if err != nil {
        return []domain.User{}, errors.New("Failed to select rows")
    }

    for rows.Next() {
        var user domain.User
        err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.FullName, &user.CreateDate, &user.UpdateDate)
        if err != nil {
            return []domain.User{}, errors.New("Failed to parse rows")
        }
        users = append(users, user)
    }

    return users, nil
}

func GetById(db *pgxpool.Pool, id uuid.UUID) (*domain.User, error) {
    var username, email, fullName string;
    var createDate, updateDate time.Time;

    err := db.QueryRow(
        context.Background(),
        "SELECT username, email, full_name, create_date, update_date FROM users WHERE user_id=$1",
        id,
    ).Scan(&username, &email, &fullName, &createDate, &updateDate);

    if err != nil {
        return nil, errors.New("Failed to get user")
    }

    return &domain.User{
        ID: id,
        Username: username,
        Email: email,
        FullName: fullName,
        CreateDate: createDate,
        UpdateDate: updateDate,
    }, nil
}

func Update(db *pgxpool.Pool, user domain.User) error {
    _, err := db.Exec(
        context.Background(),
        "UPDATE users SET username=$2, email=$3, full_name=$4 WHERE user_id=$1",
        user.ID,
        user.Username,
        user.Email,
        user.FullName,
    )

    if err != nil {
        return errors.New("Failed to update user")
    }

    return nil
}
