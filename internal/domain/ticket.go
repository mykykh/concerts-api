package domain

import (
    "time"

    "github.com/gofrs/uuid"
)

type Ticket struct {
    ID int64 `json:"id"`
    VerificationToken string `json:"verification-token"`

    ConcertID int64 `json:"concert-id"`
    UserID uuid.UUID `json:"user-id"`
    CreateDate time.Time `json:"create-date"`
    UpdateDate time.Time `json:"update-date"`
}
