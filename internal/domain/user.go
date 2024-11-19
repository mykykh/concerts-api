package domain

import (
    "time"

    "github.com/gofrs/uuid"

    "github.com/mykykh/concerts-api/internal/auth"
)

type User struct {
    ID uuid.UUID `json:"id"`

    Username string `json:"username"`
    Email string `json:"email"`
    FullName string `json:"full-name"`
    CreateDate time.Time `json:"create-date"`
    UpdateDate time.Time `json:"update-date"`
}

func ClaimsToUser(claims auth.Claims) User {
    return User{
        ID: claims.ID,
        Username: claims.Username,
        Email: claims.Email,
        FullName: claims.FullName,
    }
}
