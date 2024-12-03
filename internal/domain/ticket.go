package domain

import (
    "time"
    "errors"

    "github.com/gofrs/uuid"

    "github.com/mykykh/concerts-api/internal/utils"
)

type Ticket struct {
    ID int64 `json:"id"`
    VerificationToken string `json:"verification-token"`

    ConcertID int64 `json:"concert-id"`
    UserID uuid.UUID `json:"user-id"`
    CreateDate time.Time `json:"create-date"`
    UpdateDate time.Time `json:"update-date"`
}

func (ticket *Ticket) UpdateVerificationToken() error {
    newToken, err := utils.GenerateToken(16)
    if err != nil {
        return errors.New("Failed to generate new token")
    }
    ticket.VerificationToken = *newToken
    return nil
}
